<?php

namespace App\Console\Commands;

use App\Console\Concerns\ConfiguresSqliteDatabase;
use App\Jobs\Bible\ImportWebster1828DictionaryJob;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Bus;

class ImportDictionaryCommand extends Command
{
    use ConfiguresSqliteDatabase;

    protected $signature = 'dictionary:import
                            {--force : Re-import even if already completed}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Import bundled Webster 1828 dictionary data into the database';

    public function handle(): int
    {
        $this->configureSqliteDatabase($this->option('database'));

        Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: (bool) $this->option('force')));

        $this->info('Webster 1828 dictionary imported.');

        return self::SUCCESS;
    }
}
