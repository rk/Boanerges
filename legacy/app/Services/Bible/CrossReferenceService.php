<?php

namespace App\Services\Bible;

use App\Jobs\Bible\ImportCrossReferencesJob;
use App\Services\Bible\Import\OpenBibleVerseIdMapper;
use Illuminate\Support\Facades\DB;

class CrossReferenceService
{
    public function isImported(): bool
    {
        return DB::table('import_meta')->where('key', ImportCrossReferencesJob::META_KEY)->whereNotNull('completed_at')->exists();
    }

    /**
     * @return list<array{rank: int, bookId: string, bookName: string, chapter: int, verse: int, endVerse: int|null}>
     */
    public function forVerse(string $bookId, int $chapter, int $verse): array
    {
        $osisBookId = OsisBookId::normalize($bookId);

        if ($osisBookId === null) {
            return [];
        }

        $sourceId = $this->verseIdFor($osisBookId, $chapter, $verse);

        if ($sourceId === null) {
            return [];
        }

        $rows = DB::table('cross_references')
            ->where('source_verse_id', $sourceId)
            ->orderBy('rank')
            ->limit(100)
            ->get();

        $results = [];

        foreach ($rows as $row) {
            $target = OpenBibleVerseIdMapper::fromVerseId((int) $row->target_start_id);

            if ($target === null) {
                continue;
            }

            $endVerse = null;

            if ($row->target_end_id !== null && (int) $row->target_end_id !== (int) $row->target_start_id) {
                $end = OpenBibleVerseIdMapper::fromVerseId((int) $row->target_end_id);
                $endVerse = $end['verse'] ?? null;
            }

            $results[] = [
                'rank' => (int) $row->rank,
                'bookId' => $target['book_id'],
                'bookName' => OsisBookId::displayName($target['book_id']),
                'chapter' => $target['chapter'],
                'verse' => $target['verse'],
                'endVerse' => $endVerse,
            ];
        }

        return $results;
    }

    private function verseIdFor(string $bookId, int $chapter, int $verse): ?int
    {
        $id = DB::table('cross_reference_verses')
            ->where('book_id', $bookId)
            ->where('chapter', $chapter)
            ->where('verse', $verse)
            ->value('id');

        return $id !== null ? (int) $id : null;
    }
}
