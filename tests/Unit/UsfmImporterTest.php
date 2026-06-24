<?php

use App\Services\Bible\Import\UsfmImporter;
use Illuminate\Support\Facades\DB;
use Tests\TestCase;

uses(TestCase::class);

it('parses USFM book markers without regex errors', function () {
    $abbrev = 'usfmtest';
    $schema = app(\App\Services\Bible\TranslationSchemaManager::class);
    $schema->dropTables($abbrev);
    $schema->createTables($abbrev);

    $path = tempnam(sys_get_temp_dir(), 'usfm-');
    file_put_contents($path, <<<'USFM'
\id MAT
\h Matthew
\c 1
\p
\v 1 A genealogy of Jesus Christ, a descendant of David and Abraham.
\v 2 Abraham was the father of Isaac,
Isaac of Jacob,
\q2 Jacob of Judah and his brothers,
USFM);

    try {
        app(UsfmImporter::class)->importFromFile($abbrev, $path);

        $books = DB::table($schema->booksTable($abbrev))->get();
        $verses = DB::table($schema->versesTable($abbrev))->orderBy('verse')->get();

        expect($books)->toHaveCount(1)
            ->and($books[0]->osis_id)->toBe('mat')
            ->and($books[0]->name)->toBe('Matthew')
            ->and($verses)->toHaveCount(2)
            ->and($verses[0]->text)->toContain('genealogy')
            ->and($verses[1]->text)->toContain('Jacob of Judah')
            ->and($verses[1]->plain_text)->not->toContain('\\q2');
    } finally {
        @unlink($path);
        $schema->dropTables($abbrev);
    }
});

it('skips non-scripture USFM peripherals', function () {
    $abbrev = 'usfmperiph';
    $schema = app(\App\Services\Bible\TranslationSchemaManager::class);
    $schema->dropTables($abbrev);
    $schema->createTables($abbrev);

    $path = tempnam(sys_get_temp_dir(), 'usfm-');
    file_put_contents($path, <<<'USFM'
\id FRT
\mt1 Open English Bible
\c 1
\v 1 Not scripture
USFM);

    try {
        app(UsfmImporter::class)->importFromFile($abbrev, $path);

        expect(DB::table($schema->booksTable($abbrev))->count())->toBe(0);
    } finally {
        @unlink($path);
        $schema->dropTables($abbrev);
    }
});

it('uses the h marker for the book name', function () {
    $abbrev = 'usfmheading';
    $schema = app(\App\Services\Bible\TranslationSchemaManager::class);
    $schema->dropTables($abbrev);
    $schema->createTables($abbrev);

    $path = tempnam(sys_get_temp_dir(), 'usfm-');
    file_put_contents($path, <<<'USFM'
\id RUT Open English Bible
\ide UTF-8
\h Ruth
\mt2 The book of
\mt Ruth
\c 1
\v 1 In the days when the judges ruled there was a famine in the land.
USFM);

    try {
        app(UsfmImporter::class)->importFromFile($abbrev, $path);

        $book = DB::table($schema->booksTable($abbrev))->first();

        expect($book?->osis_id)->toBe('rut')
            ->and($book?->name)->toBe('Ruth');
    } finally {
        @unlink($path);
        $schema->dropTables($abbrev);
    }
});
