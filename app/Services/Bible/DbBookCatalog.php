<?php

namespace App\Services\Bible;

use Illuminate\Support\Facades\DB;

class DbBookCatalog
{
    public function __construct(
        private TranslationSchemaManager $schema,
    ) {}

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
    public function books(string $abbrev): array
    {
        $abbrev = $this->schema->validateAbbrev($abbrev);
        $table = $this->schema->booksTable($abbrev);

        return DB::table($table)
            ->orderBy('id')
            ->get()
            ->map(fn($book) => [
                'id' => $book->osis_id,
                'name' => $book->name,
                'abbrev' => strtoupper(substr($book->osis_id, 0, 3)),
                'testament' => $book->testament,
                'chapters' => (int) $book->chapters,
                'firstChapter' => 1,
                'lastChapter' => (int) $book->chapters,
            ])
            ->all();
    }

    /**
     * @return object{id: int, name: string, osis_id: string, testament: string, chapters: int}
     */
    public function findBook(string $abbrev, string $bookId): object
    {
        $abbrev = $this->schema->validateAbbrev($abbrev);
        $table = $this->schema->booksTable($abbrev);

        $book = DB::table($table)->where('osis_id', strtolower($bookId))->first();

        if ($book === null) {
            abort(404, "Book \"{$bookId}\" not found in {$abbrev}.");
        }

        return $book;
    }
}
