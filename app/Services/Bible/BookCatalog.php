<?php

namespace App\Services\Bible;

use rk\PhpSword\BookStructure;
use rk\PhpSword\SwordBible;

class BookCatalog
{
    /**
     * @return list<array{
     *     id: string,
     *     name: string,
     *     abbrev: string,
     *     testament: string,
     *     chapters: int,
     *     firstChapter: int,
     *     lastChapter: int
     * }>
     */
    public function books(SwordBible $bible): array
    {
        $books = [];

        foreach ($bible->getStructure()->getBooks() as $testament => $testamentBooks) {
            foreach ($testamentBooks as $book) {
                $books[] = $this->mapBook($book, $testament);
            }
        }

        return $books;
    }

    public function findBook(SwordBible $bible, string $bookId): BookStructure
    {
        [, $book] = $bible->getStructure()->findBook($bookId);

        return $book;
    }

    /**
     * @return array{
     *     id: string,
     *     name: string,
     *     abbrev: string,
     *     testament: string,
     *     chapters: int,
     *     firstChapter: int,
     *     lastChapter: int
     * }
     */
    private function mapBook(BookStructure $book, string $testament): array
    {
        return [
            'id' => strtolower($book->osisName),
            'name' => $book->name,
            'abbrev' => $book->preferredAbbreviation,
            'testament' => $testament,
            'chapters' => $book->numChapters,
            'firstChapter' => 1,
            'lastChapter' => $book->numChapters,
        ];
    }
}
