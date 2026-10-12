<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class ImportDictionaryCommand extends Command
{
    protected $signature = 'dictionary:import
                            {--force : Re-import even if already completed}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Import bundled Webster 1828 dictionary data into the database';

    public function handle(): int
    {
        $parameters = [
            'key' => 'dictionary',
            '--force' => $this->option('force'),
        ];

        $database = $this->option('database');

        if (is_string($database) && $database !== '') {
            $parameters['--database'] = $database;
        }

        return $this->call('content:import', $parameters);
    }
}
