<?php

namespace App\Console\Commands;

use App\Models\Translation;
use App\Services\Bible\Import\SwordConfReader;
use App\Services\Bible\TranslationCatalog;
use App\Services\Bible\TranslationMetadataSync;
use Illuminate\Console\Command;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;

class SyncTranslationMetadataCommand extends Command
{
    protected $signature = 'bible:sync-metadata
                            {abbrev? : Translation abbrev, e.g. asv}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Refresh translation metadata from SWORD conf files';

    public function handle(
        TranslationCatalog $catalog,
        TranslationMetadataSync $sync,
        SwordConfReader $swordConf,
    ): int {
        $this->configureDatabase($this->option('database'));

        $abbrev = $this->argument('abbrev');
        $translations = $this->matchingTranslations($abbrev, $swordConf);

        if ($translations->isEmpty()) {
            $this->error('No matching translations found.');

            return self::FAILURE;
        }

        foreach ($translations as $translation) {
            $entry = $catalog->all()->first(
                fn($item) => strcasecmp($item->short, $translation->abbrev) === 0,
            );

            $sync->applyFromSwordConf($translation, $entry);
            $this->info("Synced metadata for {$translation->abbrev}.");
        }

        DB::disconnect();

        return self::SUCCESS;
    }

    private function configureDatabase(?string $databaseOption): void
    {
        if ($databaseOption !== null) {
            $this->useDatabase($databaseOption);

            return;
        }

        if ($this->matchingTranslations(null, app(SwordConfReader::class))->isNotEmpty()) {
            return;
        }

        $nativePath = database_path('nativephp.sqlite');

        if (! is_file($nativePath)) {
            return;
        }

        $this->useDatabase($nativePath);

        if ($this->matchingTranslations(null, app(SwordConfReader::class))->isNotEmpty()) {
            $this->comment("Using NativePHP database: {$nativePath}");
        }
    }

    private function useDatabase(string $path): void
    {
        $resolved = str_starts_with($path, DIRECTORY_SEPARATOR)
            ? $path
            : base_path($path);

        config(['database.connections.sqlite.database' => $resolved]);
        DB::purge('sqlite');
        DB::reconnect('sqlite');
    }

    /** @return Collection<int, Translation> */
    private function matchingTranslations(?string $abbrev, SwordConfReader $swordConf): Collection
    {
        $query = Translation::query();

        if ($abbrev !== null) {
            $query->where('abbrev', strtolower($abbrev));
        }

        return $query->get()->filter(
            fn(Translation $translation) => $swordConf->read($translation->abbrev) !== null,
        );
    }
}
