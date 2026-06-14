<?php

use App\Services\Bible\BibleModuleManager;
use Illuminate\Support\Facades\Storage;

beforeEach(function (): void {
    if (! is_dir(Storage::disk('extras')->path('sword/mods.d'))) {
        $this->markTestSkipped('ASV SWORD module not installed. Run `php artisan bible:verify-asv`.');
    }
});

test('lists installed translations', function (): void {
    $response = $this->getJson(route('bible.translations.index'));

    $response->assertSuccessful();
    $response->assertJsonPath('translations.0.id', 'asv');
    $response->assertJsonPath('translations.0.abbrev', 'ASV');
    $response->assertJsonPath('translations.0.bundled', true);
});

test('lists books for a translation', function (): void {
    $response = $this->getJson(route('bible.books.index', ['translation' => 'asv']));

    $response->assertSuccessful();

    $genesis = collect($response->json('books'))->firstWhere('id', 'gen');

    expect($genesis)->not->toBeNull();
    expect($genesis['chapters'])->toBe(50);
    expect($genesis['firstChapter'])->toBe(1);
    expect($genesis['lastChapter'])->toBe(50);
});

test('loads genesis chapter one', function (): void {
    $response = $this->getJson(route('bible.chapters.show', [
        'translation' => 'asv',
        'book' => 'gen',
        'chapter' => 1,
    ]));

    $response->assertSuccessful();
    $response->assertJsonPath('chapter.book', 'Genesis');
    $response->assertJsonPath('chapter.chapter', 1);
    expect($response->json('chapter.verses'))->toHaveCount(31);
    expect(strtolower($response->json('chapter.verses.0.text')))->toContain('beginning');
});

test('returns not found for unknown translation', function (): void {
    $response = $this->getJson(route('bible.books.index', ['translation' => 'kjv']));

    $response->assertNotFound();
});

test('returns unprocessable for invalid chapter', function (): void {
    $response = $this->getJson(route('bible.chapters.show', [
        'translation' => 'asv',
        'book' => 'gen',
        'chapter' => 999,
    ]));

    $response->assertUnprocessable();
})->group('sword');

test('returns service unavailable when module files are missing', function (): void {
    $manager = Mockery::mock(BibleModuleManager::class);
    $manager->shouldReceive('installedModules')
        ->andReturn([
            ['key' => 'ASV', 'description' => 'American Standard Version', 'bundled' => true],
        ]);
    $manager->shouldReceive('open')
        ->once()
        ->andThrow(App\Exceptions\BibleModuleNotInstalledException::missing('ASV'));

    $this->instance(BibleModuleManager::class, $manager);

    $response = $this->getJson(route('bible.books.index', ['translation' => 'asv']));

    $response->assertServiceUnavailable();
});
