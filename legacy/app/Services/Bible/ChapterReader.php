<?php

namespace App\Services\Bible;

use App\Services\Bible\Markup\VerseTextFormatter;
use rk\PhpSword\SwordBible;

class ChapterReader
{
    public function __construct(
        private BookCatalog $bookCatalog,
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
    public function read(SwordBible $bible, string $bookId, int $chapter): array
    {
        $book = $this->bookCatalog->findBook($bible, $bookId);

        if ($chapter < 1 || $chapter > $book->numChapters) {
            abort(422, "Chapter {$chapter} is out of range for {$book->name}.");
        }

        $verseCount = $book->chapterLengths[$chapter - 1];
        $verses = [];

        for ($verseNumber = 1; $verseNumber <= $verseCount; $verseNumber++) {
            $raw = $bible->get(
                books: $book->name,
                chapters: $chapter,
                verses: $verseNumber,
                clean: false,
                join: '',
            );

            $text = $bible->get(
                books: $book->name,
                chapters: $chapter,
                verses: $verseNumber,
                clean: false,
                join: '',
            );

            $verse = [
                'number' => $verseNumber,
                'text' => trim($this->verseTextFormatter->format($text)),
            ];

            if ($this->isParagraphStart($raw)) {
                $verse['paragraphStart'] = true;
            }

            $verses[] = $verse;
        }

        return [
            'book' => $book->name,
            'bookAbbrev' => $book->preferredAbbreviation,
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
