<?php

namespace App\Console\Commands;

use App\Jobs\Bible\ImportWebster1828DictionaryJob;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Bus;

class ImportDictionaryCommand extends Command
{
    protected $signature = 'dictionary:import {--force : Re-import even if already completed}';

    protected $description = 'Import bundled Webster 1828 dictionary data into the database';

    public function handle(): int
    {
        Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: (bool) $this->option('force')));

        $this->info('Webster 1828 dictionary imported.');

        return self::SUCCESS;
    }
}
