<?php

namespace App\Console\Commands;

use App\Console\Concerns\ConfiguresSqliteDatabase;
use App\Models\Translation;
use App\Services\Bible\Import\SwordConfReader;
use App\Services\Bible\TranslationCatalog;
use App\Services\Bible\TranslationMetadataSync;
use Illuminate\Console\Command;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;

class SyncTranslationMetadataCommand extends Command
{
    use ConfiguresSqliteDatabase;

    protected $signature = 'bible:sync-metadata
                            {abbrev? : Translation abbrev, e.g. asv}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Refresh translation metadata from SWORD conf files';

    public function handle(
        TranslationCatalog $catalog,
        TranslationMetadataSync $sync,
        SwordConfReader $swordConf,
    ): int {
        $this->configureSqliteDatabase(
            $this->option('database'),
            fn(): bool => $this->matchingTranslations(null, $swordConf)->isNotEmpty(),
        );

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
