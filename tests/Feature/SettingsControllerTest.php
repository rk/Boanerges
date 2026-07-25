<?php

use App\Services\ReadabilitySettingsStore;
use Illuminate\Support\Facades\File;

test('readability settings can be retrieved', function () {
    $response = $this->getJson(route('settings.readability.show'));

    $response->assertSuccessful();
    $response->assertJsonPath('readability.fontSize', 18);
    $response->assertJsonPath('readability.theme', 'auto');
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

test('readability settings can be updated with auto theme', function () {
    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'auto',
        'fontFamily' => 'serif',
        'justifyText' => false,
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('readability.theme', 'auto');
    expect(app(ReadabilitySettingsStore::class)->get()['theme'])->toBe('auto');
});

test('readability settings can be updated with caramellatte theme', function () {
    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'caramellatte',
        'fontFamily' => 'serif',
        'justifyText' => false,
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('readability.theme', 'caramellatte');
    expect(app(ReadabilitySettingsStore::class)->get()['theme'])->toBe('caramellatte');
});

test('readability settings accept additional daisyui themes', function () {
    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'nord',
        'fontFamily' => 'serif',
        'justifyText' => true,
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('readability.theme', 'nord');

    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'dracula',
        'fontFamily' => 'serif',
        'justifyText' => true,
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('readability.theme', 'dracula');
});

test('stored sepia theme is migrated to caramellatte on read', function () {
    File::ensureDirectoryExists(dirname(storage_path('app/private/settings.json')));
    File::put(
        storage_path('app/private/settings.json'),
        json_encode([
            'readability' => [
                'theme' => 'sepia',
            ],
        ], JSON_THROW_ON_ERROR),
    );

    $response = $this->getJson(route('settings.readability.show'));

    $response->assertSuccessful();
    $response->assertJsonPath('readability.theme', 'caramellatte');

    $this->get(route('home'))
        ->assertOk()
        ->assertInertia(
            fn($page) => $page->where('readability.theme', 'caramellatte'),
        );
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

    $response = $this->patchJson(route('settings.readability.update'), [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'sepia',
        'fontFamily' => 'serif',
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
