<?php

use App\Services\StudySettingsStore;
use Illuminate\Support\Facades\File;

test('study settings can be retrieved', function () {
    $response = $this->getJson(route('settings.study.show'));

    $response->assertSuccessful();
    $response->assertJsonPath('study.activeView', 'bible');
    $response->assertJsonPath('study.bookId', 'gen');
    $response->assertJsonPath('study.chapter', 15);
});

test('study settings can be updated', function () {
    $response = $this->patchJson(route('settings.study.update'), [
        'activeView' => 'comparison',
        'bookId' => 'mat',
        'chapter' => 5,
        'translationId' => 'asv',
        'translationBId' => 'web',
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('study.activeView', 'comparison');
    $response->assertJsonPath('study.bookId', 'mat');
    $response->assertJsonPath('study.chapter', 5);

    expect(app(StudySettingsStore::class)->get()['activeView'])->toBe('comparison');
    expect(app(StudySettingsStore::class)->get()['bookId'])->toBe('mat');
});

test('study settings are shared with inertia pages', function () {
    app(StudySettingsStore::class)->update([
        'activeView' => 'scribe',
        'bookId' => 'jhn',
        'chapter' => 3,
        'translationId' => 'kjv',
        'translationBId' => 'asv',
    ]);

    $this->get(route('home'))
        ->assertOk()
        ->assertInertia(
            fn($page) => $page
            ->where('study.activeView', 'scribe')
            ->where('study.bookId', 'jhn')
            ->where('study.chapter', 3),
        );
});

test('study settings validation rejects invalid values', function () {
    $this->patchJson(route('settings.study.update'), [
        'activeView' => 'invalid',
        'bookId' => 'gen',
        'chapter' => 15,
        'translationId' => 'kjv',
        'translationBId' => 'asv',
    ])->assertUnprocessable();
});

afterEach(function () {
    $path = storage_path('app/private/settings.json');

    if (File::exists($path)) {
        File::delete($path);
    }
});
