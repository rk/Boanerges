<?php

use App\Console\Concerns\ConfiguresSqliteDatabase;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Tests\TestCase;

uses(TestCase::class);

class ConfiguresSqliteDatabaseTestCommand extends Command
{
    use ConfiguresSqliteDatabase;

    protected $signature = 'test:configure-sqlite-database';

    public function configureDatabaseForTest(?string $databaseOption, ?callable $hasDataInDefaultDatabase = null): void
    {
        $this->configureSqliteDatabase($databaseOption, $hasDataInDefaultDatabase);
    }

    public function handle(): int
    {
        return self::SUCCESS;
    }
}

test('configure sqlite database switches to explicit path', function (): void {
    $databasePath = tempnam(sys_get_temp_dir(), 'boanerges-sqlite-');
    $defaultDatabase = (string) config('database.connections.sqlite.database');
    $command = app(ConfiguresSqliteDatabaseTestCommand::class);

    try {
        $command->configureDatabaseForTest($databasePath);

        expect((string) config('database.connections.sqlite.database'))->toBe($databasePath);
    } finally {
        config(['database.connections.sqlite.database' => $defaultDatabase]);
        DB::purge('sqlite');
        DB::reconnect('sqlite');

        if ($databasePath !== false && is_file($databasePath)) {
            @unlink($databasePath);
        }
    }
});

test('configure sqlite database keeps default when callback reports data', function (): void {
    $defaultDatabase = (string) config('database.connections.sqlite.database');
    $command = app(ConfiguresSqliteDatabaseTestCommand::class);

    $command->configureDatabaseForTest(null, fn(): bool => true);

    expect((string) config('database.connections.sqlite.database'))->toBe($defaultDatabase);
});
