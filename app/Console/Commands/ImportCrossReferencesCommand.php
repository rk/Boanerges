<?php

namespace App\Console\Commands;

use App\Jobs\Bible\ImportCrossReferencesJob;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Bus;
use Illuminate\Support\Facades\DB;

class ImportCrossReferencesCommand extends Command
{
    protected $signature = 'bible:import-cross-references
                            {--force : Re-import even if already completed}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Import bundled openbible.info cross references';

    public function handle(): int
    {
        $this->configureDatabase($this->option('database'));

        Bus::dispatchSync(new ImportCrossReferencesJob(force: (bool) $this->option('force')));

        $this->info('Cross references imported.');

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
