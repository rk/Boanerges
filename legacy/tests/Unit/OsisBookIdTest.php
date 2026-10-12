<?php

use App\Services\Bible\DbBookCatalog;
use App\Services\Bible\OsisBookId;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Support\Facades\DB;
use Tests\TestCase;

uses(TestCase::class);

it('normalizes sword and readable book aliases to canonical ids', function () {
    expect(OsisBookId::normalize('matt'))->toBe('mat')
        ->and(OsisBookId::normalize('mark'))->toBe('mrk')
        ->and(OsisBookId::normalize('Matthew'))->toBe('mat')
        ->and(OsisBookId::normalize('1Sam'))->toBe('1sa')
        ->and(OsisBookId::normalize('gen'))->toBe('gen')
        ->and(OsisBookId::normalize('ezra'))->toBe('ezr')
        ->and(OsisBookId::normalize('zech'))->toBe('zec')
        ->and(OsisBookId::normalize('acts'))->toBe('act')
        ->and(OsisBookId::normalize('phlm'))->toBe('phm');
});

it('returns english display names for canonical ids', function () {
    expect(OsisBookId::displayName('ezr'))->toBe('Ezra')
        ->and(OsisBookId::displayName('zec'))->toBe('Zechariah')
        ->and(OsisBookId::displayName('mrk'))->toBe('Mark');
});

it('finds books by canonical id when translation stores sword osis ids', function () {
    $abbrev = 'osistest';
    $schema = app(TranslationSchemaManager::class);
    $schema->dropTables($abbrev);
    $schema->createTables($abbrev);

    DB::table($schema->booksTable($abbrev))->insert([
        'name' => 'Matthew',
        'osis_id' => 'matt',
        'testament' => 'nt',
        'chapters' => 28,
    ]);

    try {
        $catalog = app(DbBookCatalog::class);

        expect($catalog->findBook($abbrev, 'mat')->name)->toBe('Matthew')
            ->and($catalog->books($abbrev)[0]['id'])->toBe('mat');
    } finally {
        $schema->dropTables($abbrev);
    }
});

it('lists lookup values for legacy sword ids', function () {
    expect(OsisBookId::lookupValues('mat'))->toContain('matt', 'mat')
        ->and(OsisBookId::lookupValues('mrk'))->toContain('mark', 'mrk');
});
