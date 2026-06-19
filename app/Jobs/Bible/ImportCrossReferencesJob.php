<?php

namespace App\Jobs\Bible;

use App\Events\CrossRefImportProgress;
use App\Services\Bible\Import\OpenBibleVerseIdMapper;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Storage;

class ImportCrossReferencesJob implements ShouldQueue
{
    use Queueable;

    public int $timeout = 600;

    public function handle(): void
    {
        if (DB::table('import_meta')->where('key', 'cross_references')->whereNotNull('completed_at')->exists()) {
            return;
        }

        event(new CrossRefImportProgress('building_verse_index', 10));

        DB::table('cross_reference_verses')->delete();

        $batch = [];

        foreach (OpenBibleVerseIdMapper::all() as $id => $ref) {
            $batch[] = [
                'id' => $id,
                'book_id' => $ref['book_id'],
                'chapter' => $ref['chapter'],
                'verse' => $ref['verse'],
            ];

            if (count($batch) >= 500) {
                DB::table('cross_reference_verses')->insert($batch);
                $batch = [];
            }
        }

        if ($batch !== []) {
            DB::table('cross_reference_verses')->insert($batch);
        }

        event(new CrossRefImportProgress('importing_references', 40));

        $path = Storage::disk('extras')->path('cross-references/cross_references.txt');

        if (! is_file($path)) {
            throw new \RuntimeException('Bundled cross-reference data not found.');
        }

        DB::table('cross_references')->delete();

        $handle = fopen($path, 'r');

        if ($handle === false) {
            throw new \RuntimeException('Failed to open cross-reference data.');
        }

        $refs = [];
        $lineCount = 0;

        while (($line = fgets($handle)) !== false) {
            $parts = explode("\t", trim($line));

            if (count($parts) < 3) {
                continue;
            }

            $refs[] = [
                'source_verse_id' => (int) $parts[0],
                'rank' => (int) $parts[1],
                'target_start_id' => (int) $parts[2],
                'target_end_id' => isset($parts[3]) && $parts[3] !== '' ? (int) $parts[3] : null,
            ];

            if (count($refs) >= 1000) {
                DB::table('cross_references')->insert($refs);
                $refs = [];
            }

            $lineCount++;

            if ($lineCount % 10000 === 0) {
                event(new CrossRefImportProgress('importing_references', min(90, 40 + (int) ($lineCount / 4000))));
            }
        }

        fclose($handle);

        if ($refs !== []) {
            DB::table('cross_references')->insert($refs);
        }

        DB::table('import_meta')->updateOrInsert(
            ['key' => 'cross_references'],
            ['completed_at' => now()],
        );

        event(new CrossRefImportProgress('complete', 100));
    }
}
