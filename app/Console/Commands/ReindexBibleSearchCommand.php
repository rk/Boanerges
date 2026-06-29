<?php

namespace App\Console\Commands;

use App\Enums\TranslationInstallStatus;
use App\Models\Translation;
use App\Services\Bible\Markup\VerseTextFormatter;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Console\Command;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\DB;

class ReindexBibleSearchCommand extends Command
{
    protected $signature = 'bible:reindex-search
                            {abbrev? : Translation abbrev, e.g. asv}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Backfill plain-text verse columns and rebuild FTS indexes';

    public function handle(
        TranslationSchemaManager $schema,
        VerseTextFormatter $formatter,
    ): int {
        $this->configureDatabase($this->option('database'));

        $abbrev = $this->argument('abbrev');
        $translations = $this->readyTranslations($abbrev);

        if ($translations->isEmpty()) {
            $this->reportMissingTranslations($abbrev);

            return self::FAILURE;
        }

        foreach ($translations as $translation) {
            $this->info("Reindexing {$translation->abbrev}…");
            $schema->rebuildFtsIndex($translation->abbrev, $formatter);
        }

        $this->info('Done.');

        return self::SUCCESS;
    }

    private function configureDatabase(?string $databaseOption): void
    {
        if ($databaseOption !== null) {
            $this->useDatabase($databaseOption);

            return;
        }

        if ($this->readyTranslations(null)->isNotEmpty()) {
            return;
        }

        $nativePath = database_path('nativephp.sqlite');

        if (! is_file($nativePath)) {
            return;
        }

        $this->useDatabase($nativePath);

        if ($this->readyTranslations(null)->isNotEmpty()) {
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
    private function readyTranslations(?string $abbrev): Collection
    {
        $query = Translation::query()->where('install_status', TranslationInstallStatus::Ready);

        if ($abbrev !== null) {
            $query->where('abbrev', strtolower($abbrev));
        }

        return $query->get();
    }

    private function reportMissingTranslations(?string $abbrev): void
    {
        $database = (string) config('database.connections.sqlite.database');

        $this->error("No ready translations found in {$database}.");

        $rows = Translation::query()->get(['abbrev', 'install_status']);

        if ($rows->isEmpty()) {
            $this->line('This database has no translation rows.');

            if (is_file(database_path('nativephp.sqlite')) && $database !== database_path('nativephp.sqlite')) {
                $this->line('NativePHP stores app data in database/nativephp.sqlite. Try:');
                $this->line('  php artisan bible:reindex-search' . ($abbrev ? " {$abbrev}" : '') . ' --database=database/nativephp.sqlite');
            }

            return;
        }

        $this->line('Translations in this database:');

        foreach ($rows as $row) {
            $this->line("  {$row->abbrev}: {$row->installStatusValue()}");
        }

        if ($abbrev !== null) {
            $this->line("No ready translation matching abbrev \"{$abbrev}\".");
        }
    }
}
