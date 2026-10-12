<?php

namespace App\Services\Bible\Import;

use App\Services\Bible\Markup\VerseTextFormatter;
use App\Services\Bible\OsisBookId;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Support\Facades\DB;

class AccordanceImporter
{
    public function __construct(
        private TranslationSchemaManager $schema,
        private VerseTextFormatter $verseTextFormatter,
    ) {}

    public function importFromFile(string $abbrev, string $path): void
    {
        $abbrev = $this->schema->validateAbbrev($abbrev);
        $booksTable = $this->schema->booksTable($abbrev);
        $versesTable = $this->schema->versesTable($abbrev);

        DB::table($booksTable)->delete();
        DB::table($versesTable)->delete();

        $bookDbIds = [];
        $bookMaxChapter = [];
        $currentBook = null;

        $handle = fopen($path, 'rb');

        if ($handle === false) {
            throw new \RuntimeException('Failed to open Accordance file.');
        }

        while (($line = fgets($handle)) !== false) {
            $line = trim($line);

            if ($line === '' || str_starts_with($line, '#')) {
                continue;
            }

            if (preg_match('/^([A-Za-z0-9 ]+)\s+(\d+):(\d+)\s+(.*)$/', $line, $matches)) {
                $bookName = trim($matches[1]);
                $chapter = (int) $matches[2];
                $verse = (int) $matches[3];
                $text = trim($matches[4]);

                if ($currentBook !== $bookName) {
                    $currentBook = $bookName;
                    $osisId = $this->bookNameToOsis($bookName);
                    $bookDbIds[$bookName] = DB::table($booksTable)->insertGetId([
                        'name' => $bookName,
                        'osis_id' => $osisId,
                        'testament' => str_starts_with(strtolower($osisId), 'mat') || in_array(strtolower($osisId), ['mat', 'mrk', 'luk', 'jhn'], true) ? 'nt' : 'ot',
                        'chapters' => 0,
                    ]);
                    $bookMaxChapter[$bookName] = 0;
                }

                $bookMaxChapter[$bookName] = max($bookMaxChapter[$bookName], $chapter);
                DB::table($booksTable)->where('id', $bookDbIds[$bookName])->update([
                    'chapters' => $bookMaxChapter[$bookName],
                ]);

                DB::table($versesTable)->insert([
                    'book_id' => $bookDbIds[$bookName],
                    'chapter' => $chapter,
                    'verse' => $verse,
                    'text' => $text,
                    'plain_text' => $this->verseTextFormatter->toPlainText($text),
                ]);
            }
        }

        fclose($handle);
    }

    private function bookNameToOsis(string $name): string
    {
        return OsisBookId::normalize($name) ?? strtolower(substr(preg_replace('/[^a-z]/i', '', $name), 0, 3));
    }
}
