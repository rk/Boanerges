<?php

return [
    'bundled_sword_path' => 'sword',
    'local_sword_path' => 'modules/bible',
    'catalog_path' => 'translations.json',

    'bundled_modules' => [
        'ASV',
    ],

    // When true in testing, EnsureBundledData may run real SWORD imports (slow).
    // Feature tests use BundledAsvTestSeeder instead; see Phase 2 SQL fixture plan there.
    'seed_bundled_in_tests' => env('BOANERGES_SEED_BUNDLED_IN_TESTS', false),

    'readability' => [
        'fontSize' => 18,
        'lineHeight' => 1.7,
        'theme' => 'light',
        'fontFamily' => 'serif',
    ],

    'study' => [
        'activeView' => 'bible',
        'bookId' => 'gen',
        'chapter' => 15,
        'translationId' => 'asv',
        'translationBId' => 'asv',
    ],
];
