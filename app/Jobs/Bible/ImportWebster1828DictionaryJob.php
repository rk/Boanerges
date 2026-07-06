<?php

namespace App\Jobs\Bible;

use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Storage;

class ImportWebster1828DictionaryJob implements ShouldQueue
{
    use Queueable;

    public const META_KEY = 'webster1828_dictionary_v1';

    public int $timeout = 600;

    public function __construct(
        private bool $force = false,
    ) {}

    public function handle(): void
    {
        if (
            ! $this->force
            && DB::table('import_meta')->where('key', self::META_KEY)->whereNotNull('completed_at')->exists()
        ) {
            return;
        }

        $path = $this->dictionaryPath();

        if (! is_file($path)) {
            throw new \RuntimeException('Bundled Webster 1828 dictionary data not found.');
        }

        $entries = json_decode((string) file_get_contents($path), true);

        if (! is_array($entries)) {
            throw new \RuntimeException('Failed to decode Webster 1828 dictionary data.');
        }

        DB::table('dictionary_entries')->delete();

        $batch = [];

        foreach ($entries as $entry) {
            if (! is_array($entry)) {
                continue;
            }

            $word = trim((string) ($entry['word'] ?? ''));

            if ($word === '') {
                continue;
            }

            $definitions = $entry['definitions'] ?? [];

            if (! is_array($definitions)) {
                continue;
            }

            $batch[] = [
                'word' => $word,
                'word_key' => mb_strtolower($word),
                'part_of_speech' => isset($entry['pos']) ? (string) $entry['pos'] : null,
                'definitions' => json_encode(array_values(array_map(
                    static fn($definition) => (string) $definition,
                    $definitions,
                ))),
            ];

            if (count($batch) >= 1000) {
                DB::table('dictionary_entries')->insert($batch);
                $batch = [];
            }
        }

        if ($batch !== []) {
            DB::table('dictionary_entries')->insert($batch);
        }

        DB::table('import_meta')->updateOrInsert(
            ['key' => self::META_KEY],
            ['completed_at' => now()],
        );
    }

    private function dictionaryPath(): string
    {
        if (app()->environment('testing')) {
            return base_path('tests/fixtures/websters1828/dictionary.json');
        }

        return Storage::disk('extras')->path((string) config('boanerges.dictionary_path'));
    }
}
