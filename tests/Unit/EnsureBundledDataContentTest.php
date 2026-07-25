<?php

use App\Jobs\Bible\ImportCrossReferencesJob;
use App\Listeners\EnsureBundledData;
use App\Support\BundledContentRegistry;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Support\Facades\Bus;

uses(Tests\TestCase::class, RefreshDatabase::class);

test('ensure bundled data dispatches pending content when seed_bundled_in_tests is enabled', function () {
    config([
        'boanerges.seed_bundled_in_tests' => true,
        'boanerges.bundled_modules' => [],
        'queue.default' => 'sync',
    ]);

    Bus::fake();

    $crossRefs = Mockery::mock(\App\Contracts\BundledContentProvider::class);
    $crossRefs->shouldReceive('key')->andReturn('cross-references');
    $crossRefs->shouldReceive('isImported')->andReturn(false);
    $crossRefs->shouldReceive('requiresReadyTranslation')->andReturn(false);
    $crossRefs->shouldReceive('importJob')->once()->andReturn(new ImportCrossReferencesJob(force: true));

    $registry = new BundledContentRegistry();
    $registry->register($crossRefs);

    $installer = Mockery::mock(\App\Services\Bible\TranslationInstaller::class);
    $installer->shouldReceive('installBundled')->never();

    (new EnsureBundledData($installer, $registry))->handle();

    Bus::assertDispatched(ImportCrossReferencesJob::class);
});
