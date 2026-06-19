<?php

namespace App\Services\Bible;

use App\Services\Bible\Markup\VerseTextFormatter;
use Illuminate\Support\Facades\DB;

class DbChapterReader
{
    public function __construct(
        private DbBookCatalog $bookCatalog,
        private TranslationSchemaManager $schema,
        private VerseTextFormatter $verseTextFormatter,
    ) {}

    /**
     * @return array{
     *     book: string,
     *     bookAbbrev: string,
     *     chapter: int,
     *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
     * }
     */
    public function read(string $abbrev, string $bookId, int $chapter): array
    {
        $book = $this->bookCatalog->findBook($abbrev, $bookId);

        if ($chapter < 1 || $chapter > $book->chapters) {
            abort(422, "Chapter {$chapter} is out of range for {$book->name}.");
        }

        $versesTable = $this->schema->versesTable($abbrev);

        $rows = DB::table($versesTable)
            ->where('book_id', $book->id)
            ->where('chapter', $chapter)
            ->orderBy('verse')
            ->get();

        $verses = [];

        foreach ($rows as $row) {
            $raw = (string) $row->text;
            $verse = [
                'number' => (int) $row->verse,
                'text' => trim($this->verseTextFormatter->format($raw)),
            ];

            if ($this->isParagraphStart($raw)) {
                $verse['paragraphStart'] = true;
            }

            $verses[] = $verse;
        }

        return [
            'book' => $book->name,
            'bookAbbrev' => strtoupper(substr($book->osis_id, 0, 3)),
            'chapter' => $chapter,
            'verses' => $verses,
        ];
    }

    private function isParagraphStart(string $rawText): bool
    {
        $trimmed = ltrim($rawText);

        return (bool) preg_match('/<(p|div)\b/i', $trimmed)
            || str_contains($trimmed, 'type="x-p"');
    }
}
