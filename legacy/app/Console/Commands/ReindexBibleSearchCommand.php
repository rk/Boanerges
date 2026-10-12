<?php

namespace App\Console\Commands;

use App\Console\Concerns\ConfiguresSqliteDatabase;
use App\Enums\TranslationInstallStatus;
use App\Models\Translation;
use App\Services\Bible\Markup\VerseTextFormatter;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Console\Command;
use Illuminate\Support\Collection;

class ReindexBibleSearchCommand extends Command
{
    use ConfiguresSqliteDatabase;

    protected $signature = 'bible:reindex-search
                            {abbrev? : Translation abbrev, e.g. asv}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Backfill plain-text verse columns and rebuild FTS indexes';

    public function handle(
        TranslationSchemaManager $schema,
        VerseTextFormatter $formatter,
    ): int {
        $this->configureSqliteDatabase(
            $this->option('database'),
            fn(): bool => $this->readyTranslations(null)->isNotEmpty(),
        );

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
