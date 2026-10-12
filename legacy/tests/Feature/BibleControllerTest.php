<?php

use App\Enums\TranslationInstallStatus;
use App\Models\Translation;

beforeEach(function (): void {
    seedBundledAsvForTests();
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

test('search finds genesis reference text', function (): void {
    $response = $this->getJson(route('bible.search', ['q' => 'beginning', 'translation' => 'asv']));

    $response->assertSuccessful();
    expect(collect($response->json('results'))->firstWhere('bookId', 'gen'))->not->toBeNull();
});

test('search supports pagination metadata', function (): void {
    $response = $this->getJson(route('bible.search', [
        'q' => 'beginning',
        'translation' => 'asv',
        'limit' => 5,
        'offset' => 0,
    ]));

    $response->assertSuccessful();
    $response->assertJsonStructure(['results', 'total', 'hasMore']);
    expect($response->json('results'))->not->toBeEmpty();
});

test('returns not found when translation is not ready', function (): void {
    Translation::query()->where('abbrev', 'asv')->update([
        'install_status' => TranslationInstallStatus::Importing,
    ]);

    $response = $this->getJson(route('bible.books.index', ['translation' => 'asv']));

    $response->assertNotFound();
});
