<?php

use App\Enums\StudyColumnType;

uses(Tests\TestCase::class);

test('frontend column catalog matches StudyColumnType enum', function () {
    $catalogPath = base_path('resources/js/lib/columns/catalog.ts');
    $contents = file_get_contents($catalogPath);

    expect($contents)->not->toBeFalse();

    preg_match_all("/type:\s*'([^']+)'/", $contents, $matches);

    $frontendTypes = array_values(array_unique($matches[1] ?? []));
    sort($frontendTypes);

    $phpTypes = array_map(
        static fn(StudyColumnType $type): string => $type->value,
        StudyColumnType::cases(),
    );
    sort($phpTypes);

    expect($frontendTypes)->toBe($phpTypes);
});

test('frontend menu ids match StudyColumnType menuId values', function () {
    $catalogPath = base_path('resources/js/lib/columns/catalog.ts');
    $contents = (string) file_get_contents($catalogPath);

    preg_match_all("/menuId:\s*'([^']+)'/", $contents, $matches);
    $frontendMenuIds = $matches[1] ?? [];
    sort($frontendMenuIds);

    $phpMenuIds = [];

    foreach (StudyColumnType::cases() as $type) {
        $menuId = $type->menuId();

        if ($menuId !== null) {
            $phpMenuIds[] = $menuId;
        }
    }

    sort($phpMenuIds);

    expect($frontendMenuIds)->toBe($phpMenuIds);
});

test('study menu column order matches pre-refactor Study menu order', function () {
    $menuTypes = [];

    foreach (StudyColumnType::cases() as $type) {
        if ($type->menuId() !== null) {
            $menuTypes[] = $type->value;
        }
    }

    expect($menuTypes)->toBe([
        'search',
        'cross-references',
        'comparison',
        'verse-list',
        'dictionary',
    ]);

    expect(StudyColumnType::CrossReferences->menuLabel())->toBe('Cross-References');
});
