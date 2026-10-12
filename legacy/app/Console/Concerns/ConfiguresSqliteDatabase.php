<?php

namespace App\Console\Concerns;

use Illuminate\Support\Facades\DB;

trait ConfiguresSqliteDatabase
{
    protected function configureSqliteDatabase(?string $databaseOption, ?callable $hasDataInDefaultDatabase = null): void
    {
        if (is_string($databaseOption) && $databaseOption !== '') {
            $this->useSqliteDatabase($databaseOption);

            return;
        }

        if ($hasDataInDefaultDatabase !== null && $hasDataInDefaultDatabase()) {
            return;
        }

        // Keep the in-memory test connection; NativePHP DB would poison RefreshDatabase.
        if (app()->environment('testing')) {
            return;
        }

        $nativePath = database_path('nativephp.sqlite');

        if (! is_file($nativePath)) {
            return;
        }

        $this->useSqliteDatabase($nativePath);

        if ($hasDataInDefaultDatabase === null || $hasDataInDefaultDatabase()) {
            $this->comment("Using NativePHP database: {$nativePath}");
        }
    }

    protected function useSqliteDatabase(string $path): void
    {
        $resolved = $path === ':memory:' || str_starts_with($path, DIRECTORY_SEPARATOR)
            ? $path
            : base_path($path);

        config(['database.connections.sqlite.database' => $resolved]);
        DB::purge('sqlite');
        DB::reconnect('sqlite');
    }
}
