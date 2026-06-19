<?php

use App\Services\Bible\BibleModuleManager;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Storage;

beforeEach(function (): void {
    if (! is_dir(Storage::disk('extras')->path('sword/mods.d'))) {
        $this->markTestSkipped('ASV SWORD module not installed.');
    }
});

afterEach(function (): void {
    $translation = \App\Models\Translation::query()->where('abbrev', 'kjv')->first();

    if ($translation !== null) {
        app(\App\Services\Bible\TranslationSchemaManager::class)->dropTables('kjv');
        $translation->delete();
    }

    app(BibleModuleManager::class)->clearCache();
});

test('catalog lists english translations with install status', function (): void {
    $response = $this->getJson(route('bible.translations.catalog'));

    $response->assertSuccessful();

    $asv = collect($response->json('translations'))->firstWhere('module', 'ASV');

    expect($asv)->not->toBeNull();
    expect($asv['installed'])->toBeTrue();
    expect($asv['bundled'])->toBeTrue();
    expect(collect($response->json('translations'))->count())->toBeGreaterThan(50);
});

test('installed translations are discovered from database', function (): void {
    $response = $this->getJson(route('bible.translations.index'));

    $response->assertSuccessful();
    $response->assertJsonPath('translations.0.id', 'asv');
    $response->assertJsonPath('translations.0.bundled', true);
});

test('bundled asv cannot be installed via api', function (): void {
    $response = $this->postJson(route('bible.translations.install', ['module' => 'ASV']));

    $response->assertUnprocessable();
});

test('bundled asv cannot be uninstalled via api', function (): void {
    $response = $this->deleteJson(route('bible.translations.uninstall', ['module' => 'ASV']));

    $response->assertUnprocessable();
});

test('translation can be installed and uninstalled', function (): void {
    $fixture = '/Users/robert/Code/php-sword/tests/resources/KJV.zip';

    if (! is_file($fixture)) {
        $this->markTestSkipped('KJV fixture zip not available.');
    }

    $catalog = app(\App\Services\Bible\TranslationCatalog::class)->find('KJV');

    Http::fake([
        $catalog->url => Http::response(file_get_contents($fixture)),
    ]);

    $install = $this->postJson(route('bible.translations.install', ['module' => 'KJV']));

    $install->assertAccepted();
    $install->assertJsonPath('translation.id', 'kjv');

    $this->getJson(route('bible.translations.index'))
        ->assertJsonFragment(['id' => 'kjv']);

    $this->deleteJson(route('bible.translations.uninstall', ['module' => 'KJV']))
        ->assertSuccessful();

    $this->getJson(route('bible.translations.index'))
        ->assertJsonMissing(['id' => 'kjv']);
})->group('sword');

test('installing an already installed translation returns conflict', function (): void {
    $fixture = '/Users/robert/Code/php-sword/tests/resources/KJV.zip';

    if (! is_file($fixture)) {
        $this->markTestSkipped('KJV fixture zip not available.');
    }

    $catalog = app(\App\Services\Bible\TranslationCatalog::class)->find('KJV');

    Http::fake([
        $catalog->url => Http::response(file_get_contents($fixture)),
    ]);

    $this->postJson(route('bible.translations.install', ['module' => 'KJV']))->assertAccepted();

    $this->postJson(route('bible.translations.install', ['module' => 'KJV']))
        ->assertConflict();
})->group('sword');
