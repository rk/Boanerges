<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;

class ImportCrossReferencesCommand extends Command
{
    protected $signature = 'bible:import-cross-references
                            {--force : Re-import even if already completed}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Import bundled openbible.info cross references';

    public function handle(): int
    {
        $parameters = [
            'key' => 'cross-references',
            '--force' => $this->option('force'),
        ];

        $database = $this->option('database');

        if (is_string($database) && $database !== '') {
            $parameters['--database'] = $database;
        }

        return $this->call('content:import', $parameters);
    }
}
