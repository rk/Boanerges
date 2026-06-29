<?php

use App\Services\ReadabilitySettingsStore;
use Illuminate\Support\Facades\File;

test('readability settings can be retrieved', function () {
    $response = $this->getJson(route('settings.readability.show'));

    $response->assertSuccessful();
    $response->assertJsonPath('readability.fontSize', 18);
    $response->assertJsonPath('readability.fontFamily', 'serif');
});

test('readability settings can be updated', function () {
    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 20,
        'lineHeight' => 1.8,
        'theme' => 'dark',
        'fontFamily' => 'sans-serif',
        'justifyText' => true,
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('readability.fontSize', 20);
    $response->assertJsonPath('readability.theme', 'dark');
    $response->assertJsonPath('readability.fontFamily', 'sans-serif');

    expect(app(ReadabilitySettingsStore::class)->get()['fontSize'])->toBe(20);
    expect(app(ReadabilitySettingsStore::class)->get()['theme'])->toBe('dark');
    expect(app(ReadabilitySettingsStore::class)->get()['fontFamily'])->toBe('sans-serif');
});

test('readability settings are shared with inertia pages', function () {
    app(ReadabilitySettingsStore::class)->update([
        'fontSize' => 22,
        'lineHeight' => 1.9,
        'theme' => 'dark',
        'fontFamily' => 'sans-serif',
    ]);

    $this->get(route('home'))
        ->assertOk()
        ->assertInertia(
            fn($page) => $page
            ->where('readability.fontSize', 22)
            ->where('readability.fontFamily', 'sans-serif'),
        );
});

test('readability settings can be updated with sepia theme', function () {
    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'sepia',
        'fontFamily' => 'serif',
        'justifyText' => false,
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('readability.theme', 'sepia');
    expect(app(ReadabilitySettingsStore::class)->get()['theme'])->toBe('sepia');
});

test('readability settings validation rejects invalid values', function () {
    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 10,
        'lineHeight' => 1.8,
        'theme' => 'dark',
        'fontFamily' => 'sans-serif',
        'justifyText' => true,
    ]);

    $response->assertUnprocessable();
});

afterEach(function () {
    $path = storage_path('app/private/settings.json');

    if (File::exists($path)) {
        File::delete($path);
    }
});
