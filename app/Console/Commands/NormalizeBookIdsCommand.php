<?php

namespace App\Console\Commands;

use App\Enums\TranslationInstallStatus;
use App\Models\Translation;
use App\Services\Bible\OsisBookId;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;

class NormalizeBookIdsCommand extends Command
{
    protected $signature = 'bible:normalize-book-ids
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Rewrite per-translation book osis_id values to canonical IDs';

    public function handle(TranslationSchemaManager $schema): int
    {
        $this->configureDatabase($this->option('database'));

        $updated = 0;

        foreach (Translation::query()->where('install_status', TranslationInstallStatus::Ready)->get() as $translation) {
            $table = $schema->booksTable($translation->abbrev);

            foreach (DB::table($table)->get(['id', 'osis_id']) as $book) {
                $canonical = OsisBookId::normalize((string) $book->osis_id);

                if ($canonical === null || $canonical === $book->osis_id) {
                    continue;
                }

                DB::table($table)->where('id', $book->id)->update(['osis_id' => $canonical]);
                $updated++;
            }
        }

        $this->info("Updated {$updated} book rows.");

        return self::SUCCESS;
    }

    private function configureDatabase(?string $databaseOption): void
    {
        if ($databaseOption !== null) {
            $this->useDatabase($databaseOption);

            return;
        }

        $nativePath = database_path('nativephp.sqlite');

        if (is_file($nativePath)) {
            $this->useDatabase($nativePath);
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
}
