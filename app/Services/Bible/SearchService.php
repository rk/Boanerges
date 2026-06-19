<?php

namespace App\Services\Bible;

use App\Enums\TranslationInstallStatus;
use App\Models\Translation;
use Illuminate\Support\Facades\DB;

class SearchService
{
    public function __construct(
        private TranslationSchemaManager $schema,
    ) {}

    /**
     * @return list<array{bookId: string, chapter: int, verse: int, snippet: string, translation: string}>
     */
    public function search(string $query, ?string $translation = null, int $limit = 50): array
    {
        $translations = $translation !== null
            ? [Translation::query()->where('abbrev', strtolower($translation))->where('install_status', TranslationInstallStatus::Ready)->firstOrFail()]
            : Translation::query()->where('install_status', TranslationInstallStatus::Ready)->get()->all();

        $results = [];

        foreach ($translations as $item) {
            $fts = $this->schema->ftsTable($item->abbrev);

            if (! in_array($fts, $this->tableNames(), true)) {
                continue;
            }

            $rows = DB::select(
                "SELECT book_id, chapter, verse, snippet({$fts}, 0, '<mark>', '</mark>', '…', 32) AS snippet
                 FROM {$fts}
                 WHERE {$fts} MATCH ?
                 LIMIT ?",
                [$query, $limit],
            );

            foreach ($rows as $row) {
                $bookRow = DB::table($this->schema->booksTable($item->abbrev))
                    ->where('id', $row->book_id)
                    ->first();

                if ($bookRow === null) {
                    continue;
                }

                $results[] = [
                    'bookId' => $bookRow->osis_id,
                    'chapter' => (int) $row->chapter,
                    'verse' => (int) $row->verse,
                    'snippet' => strip_tags((string) $row->snippet, '<mark>'),
                    'translation' => $item->abbrev,
                ];
            }
        }

        return array_slice($results, 0, $limit);
    }

    /** @return list<string> */
    private function tableNames(): array
    {
        return array_column(
            DB::select("SELECT name FROM sqlite_master WHERE type='table'"),
            'name',
        );
    }
}
