<?php

namespace App\Console\Commands;

use App\Console\Concerns\ConfiguresSqliteDatabase;
use App\Jobs\Bible\ImportCrossReferencesJob;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Bus;

class ImportCrossReferencesCommand extends Command
{
    use ConfiguresSqliteDatabase;

    protected $signature = 'bible:import-cross-references
                            {--force : Re-import even if already completed}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Import bundled openbible.info cross references';

    public function handle(): int
    {
        $this->configureSqliteDatabase($this->option('database'));

        Bus::dispatchSync(new ImportCrossReferencesJob(force: (bool) $this->option('force')));

        $this->info('Cross references imported.');

        return self::SUCCESS;
    }
}
