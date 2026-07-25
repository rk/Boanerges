<?php

use App\Support\BundledContentRegistry;
use Illuminate\Foundation\Testing\RefreshDatabase;

uses(RefreshDatabase::class, Tests\TestCase::class);

test('registry exposes dictionary and cross reference providers', function () {
    $registry = app(BundledContentRegistry::class);
    $keys = array_map(static fn($provider) => $provider->key(), $registry->all());

    expect($keys)->toContain('dictionary', 'cross-references');
});

test('pending providers skip already imported content', function () {
    $registry = app(BundledContentRegistry::class);

    foreach ($registry->all() as $provider) {
        if ($provider->isImported()) {
            expect(collect($registry->pending())->pluck('key'))->not->toContain($provider->key());
        }
    }
});
