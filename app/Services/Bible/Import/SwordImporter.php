<?php

namespace App\Services\Bible\Import;

use App\Services\BatchedInsertQueue;
use App\Services\Bible\BibleModuleManager;
use App\Services\Bible\OsisBookId;
use App\Services\Bible\Markup\VerseTextFormatter;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Support\Facades\DB;

class SwordImporter
{
    use ProgressEmitter;

    public function __construct(
        private BibleModuleManager $modules,
        private TranslationSchemaManager $schema,
        private VerseTextFormatter $verseTextFormatter,
    ) {}

    public function import(string $abbrev, string $moduleKey): void
    {
        $abbrev = $this->schema->validateAbbrev($abbrev);
        $bible = $this->modules->open($moduleKey);
        $booksTable = $this->schema->booksTable($abbrev);
        $versesTable = $this->schema->versesTable($abbrev);

        DB::table($booksTable)->delete();
        DB::table($versesTable)->delete();

        $processor = new BatchedInsertQueue(
            static fn (array $rows) => DB::table($versesTable)->insert($rows),
        );

        $structure = $bible->getStructure()->getBooks();

        $maximum = 0;
        $counter = 0;

        foreach ($structure as $books) {
            $maximum += count($books);
        }

        foreach ($structure as $testament => $testamentBooks) {
            foreach ($testamentBooks as $book) {
                $bookId = DB::table($booksTable)->insertGetId([
                    'name' => $book->name,
                    'osis_id' => OsisBookId::normalize(strtolower($book->osisName)) ?? strtolower($book->osisName),
                    'testament' => strtolower($testament),
                    'chapters' => $book->numChapters,
                ]);

                for ($chapter = 1; $chapter <= $book->numChapters; $chapter++) {
                    $verseCount = $book->chapterLengths[$chapter - 1];

                    for ($verse = 1; $verse <= $verseCount; $verse++) {
                        $text = $bible->get(
                            books: $book->name,
                            chapters: $chapter,
                            verses: $verse,
                            clean: false,
                            join: '',
                        );

                        $processor->push([
                            'book_id' => $bookId,
                            'chapter' => $chapter,
                            'verse' => $verse,
                            'text' => $text,
                            'plain_text' => $this->verseTextFormatter->toPlainText($text),
                        ]);
                    }
                }

                $this->onProgress(++$counter, $maximum);
            }
        }

        $processor->done();
    }

    public function verify(string $moduleKey): void
    {
        $bible = $this->modules->open($moduleKey);

        $text = trim($bible->get(
            books: 'Genesis',
            chapters: 1,
            verses: 1,
            clean: true,
            join: '',
        ));

        $text .= trim($bible->get(
            books: 'Matthew',
            chapters: 1,
            verses: 1,
            clean: true,
            join: '',
        ));

        if ($text === '') {
            throw new \RuntimeException("Verification failed for {$moduleKey}.");
        }
    }
}
