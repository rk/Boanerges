<?php

use App\Services\StudySettingsStore;
use Illuminate\Support\Facades\File;

test('study settings can be retrieved', function () {
    $response = $this->getJson(route('settings.study.show'));

    $response->assertSuccessful();
    $response->assertJsonPath('study.columnCount', 1);
    $response->assertJsonPath('study.columns', []);
    $response->assertJsonPath('study.bookId', 'gen');
    $response->assertJsonPath('study.chapter', 15);
});

test('study settings can be updated', function () {
    $response = $this->patchJson(route('settings.study.update'), [
        'columnCount' => 2,
        'columns' => ['bible-secondary'],
        'bookId' => 'mat',
        'chapter' => 5,
        'translationId' => 'asv',
        'translationBId' => 'web',
        'translationCId' => 'asv',
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('study.columnCount', 2);
    $response->assertJsonPath('study.columns', ['bible-secondary']);
    $response->assertJsonPath('study.bookId', 'mat');
    $response->assertJsonPath('study.chapter', 5);

    expect(app(StudySettingsStore::class)->get()['columnCount'])->toBe(2);
    expect(app(StudySettingsStore::class)->get()['bookId'])->toBe('mat');
});

test('study settings are shared with inertia pages', function () {
    app(StudySettingsStore::class)->update([
        'columnCount' => 3,
        'columns' => ['scribe', 'bible-secondary'],
        'bookId' => 'jhn',
        'chapter' => 3,
        'translationId' => 'kjv',
        'translationBId' => 'asv',
        'translationCId' => 'web',
    ]);

    $this->get(route('home'))
        ->assertOk()
        ->assertInertia(
            fn($page) => $page
                ->where('study.columnCount', 3)
                ->where('study.bookId', 'jhn')
                ->where('study.chapter', 3),
        );
});

test('legacy activeView settings migrate on read', function () {
    $path = storage_path('app/private/settings.json');
    File::ensureDirectoryExists(dirname($path));
    File::put($path, json_encode([
        'study' => [
            'activeView' => 'comparison',
            'bookId' => 'mat',
            'chapter' => 5,
            'translationId' => 'asv',
            'translationBId' => 'web',
        ],
    ], JSON_THROW_ON_ERROR));

    $settings = app(StudySettingsStore::class)->get();

    expect($settings['columnCount'])->toBe(2);
    expect($settings['columns'])->toBe(['bible-secondary']);
    expect($settings['bookId'])->toBe('mat');
});

test('invalid saved column layout is sanitized on read', function () {
    app(StudySettingsStore::class)->update([
        'columnCount' => 2,
        'columns' => ['invalid-type'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ]);

    $settings = app(StudySettingsStore::class)->get();

    expect($settings['columnCount'])->toBe(2);
    expect($settings['columns'])->toBe(['bible-secondary']);
});

test('study settings accept cross-references column type', function () {
    $response = $this->patchJson(route('settings.study.update'), [
        'columnCount' => 2,
        'columns' => ['cross-references'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('study.columns', ['cross-references']);
});

test('study settings validation rejects invalid values', function () {
    $this->patchJson(route('settings.study.update'), [
        'columnCount' => 2,
        'columns' => [],
        'bookId' => 'gen',
        'chapter' => 15,
        'translationId' => 'kjv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ])->assertUnprocessable();
});

afterEach(function () {
    $path = storage_path('app/private/settings.json');

    if (File::exists($path)) {
        File::delete($path);
    }
});
