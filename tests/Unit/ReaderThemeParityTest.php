<?php

use App\Enums\ReaderTheme;

uses(Tests\TestCase::class);

test('reader themes match frontend catalog and daisyui css list', function () {
    $phpThemes = array_map(
        static fn(ReaderTheme $theme): string => $theme->value,
        ReaderTheme::cases(),
    );
    sort($phpThemes);

    $themesTs = file_get_contents(base_path('resources/js/lib/themes.ts'));
    preg_match_all("/id: '([^']+)'/", $themesTs, $themeMatches);
    $frontendThemes = array_values(array_unique($themeMatches[1] ?? []));
    sort($frontendThemes);

    expect($frontendThemes)->toBe($phpThemes);

    $css = file_get_contents(base_path('resources/css/app.css'));
    preg_match('/@plugin \'daisyui\' \{[^}]*themes:\s*([^;]+);/s', $css, $cssMatch);
    $cssThemes = preg_split('/\s*,\s*/', trim($cssMatch[1] ?? ''));
    $cssThemes = array_values(array_filter(array_map(
        static fn(string $theme): string => preg_replace('/\s+--\S+/', '', trim($theme)),
        $cssThemes,
    )));
    sort($cssThemes);

    $daisyThemes = array_values(array_filter(
        $phpThemes,
        static fn(string $theme): bool => ! in_array($theme, ['auto'], true),
    ));
    sort($daisyThemes);

    expect($cssThemes)->toBe($daisyThemes);
});
