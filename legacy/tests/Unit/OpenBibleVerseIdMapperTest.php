<?php

use App\Services\Bible\Import\OpenBibleVerseIdMapper;

it('maps openbible references to verse ids', function () {
    expect(OpenBibleVerseIdMapper::verseIdsFromReference('Gen.1.1'))->toBe([
        'start' => 1,
        'end' => null,
    ])->and(OpenBibleVerseIdMapper::verseIdsFromReference('John.1.1-John.1.3'))->toBe([
        'start' => 26103,
        'end' => 26105,
    ]);
});

it('normalizes common book aliases', function () {
    expect(OpenBibleVerseIdMapper::normalizeOsisBookId('mark'))->toBe('mrk')
        ->and(OpenBibleVerseIdMapper::normalizeOsisBookId('Gen'))->toBe('gen')
        ->and(OpenBibleVerseIdMapper::normalizeOsisBookId('genesis'))->toBe('gen');
});
