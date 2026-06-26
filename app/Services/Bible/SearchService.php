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
     * @return array{results: list<array{bookId: string, chapter: int, verse: int, snippet: string, translation: string}>, total: int, hasMore: bool}
     */
    public function search(
        string $query,
        ?string $translation = null,
        int $limit = 50,
        int $offset = 0,
    ): array {
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
                 LIMIT ? OFFSET ?",
                [$query, $limit + 1, $offset],
            );

            foreach ($rows as $row) {
                $bookRow = DB::table($this->schema->booksTable($item->abbrev))
                    ->where('id', $row->book_id)
                    ->first();

                if ($bookRow === null) {
                    continue;
                }

                $results[] = [
                    'bookId' => OsisBookId::normalize($bookRow->osis_id) ?? $bookRow->osis_id,
                    'chapter' => (int) $row->chapter,
                    'verse' => (int) $row->verse,
                    'snippet' => strip_tags((string) $row->snippet, '<mark>'),
                    'translation' => $item->abbrev,
                ];
            }
        }

        $hasMore = count($results) > $limit;
        $page = array_slice($results, 0, $limit);
        $total = $offset + count($page) + ($hasMore ? 1 : 0);

        return [
            'results' => $page,
            'total' => $total,
            'hasMore' => $hasMore,
        ];
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
