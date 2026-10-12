<?php

namespace App\Console\Commands;

use App\Console\Concerns\ConfiguresSqliteDatabase;
use App\Support\BundledContentRegistry;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Bus;

class ImportContentCommand extends Command
{
    use ConfiguresSqliteDatabase;

    protected $signature = 'content:import
                            {key? : Provider key (dictionary, cross-references)}
                            {--force : Re-import even if already completed}
                            {--database= : SQLite file path (defaults to NativePHP dev DB when present)}';

    protected $description = 'Import bundled content providers';

    public function handle(BundledContentRegistry $registry): int
    {
        $this->configureSqliteDatabase($this->option('database'));

        $key = $this->argument('key');
        $force = (bool) $this->option('force');

        $providers = $key !== null
            ? [$registry->get((string) $key)]
            : $registry->all();

        foreach ($providers as $provider) {
            Bus::dispatchSync($provider->importJob($force));
            $this->info("Imported {$provider->key()}.");
        }

        return self::SUCCESS;
    }
}
