<?php

namespace App\Services\Bible\Import;

class OpenBibleVerseIdMapper
{
    /** @var array<int, array{book_id: string, chapter: int, verse: int}>|null */
    private static ?array $lookup = null;

    /**
     * @return array{book_id: string, chapter: int, verse: int}|null
     */
    public static function fromVerseId(int $verseId): ?array
    {
        self::$lookup ??= self::buildLookup();

        return self::$lookup[$verseId] ?? null;
    }

    /** @return array<int, array{book_id: string, chapter: int, verse: int}> */
    public static function all(): array
    {
        self::$lookup ??= self::buildLookup();

        return self::$lookup;
    }

    /** @return array<int, array{book_id: string, chapter: int, verse: int}> */
    private static function buildLookup(): array
    {
        $lookup = [];
        $id = 1;

        foreach (BibleCanon::books() as $book) {
            foreach ($book['chapters'] as $chapterIndex => $verseCount) {
                $chapter = $chapterIndex + 1;

                for ($verse = 1; $verse <= $verseCount; $verse++) {
                    $lookup[$id] = [
                        'book_id' => $book['osis'],
                        'chapter' => $chapter,
                        'verse' => $verse,
                    ];
                    $id++;
                }
            }
        }

        return $lookup;
    }
}
