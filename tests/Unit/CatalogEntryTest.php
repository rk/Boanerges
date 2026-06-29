<?php

use App\Data\CatalogEntry;
use App\Enums\CatalogImportFormat;

it('maps OEB catalog entry to USFM import and markup', function () {
    $entry = CatalogEntry::fromArray([
        'short' => 'OEB',
        'name' => 'Open English Bible (US Spelling)',
        'url' => 'https://openenglishbible.org/oeb/2025.6/OEB-2025.6-US.usfm.zip',
        'import_as' => 'usfm',
        'markup_format' => 'usfm',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Usfm)
        ->and($entry->markupFormat)->toBe('usfm');
});

it('defaults USFM verse markup when only import format is given', function () {
    $entry = CatalogEntry::fromArray([
        'short' => 'OEB',
        'name' => 'Open English Bible (US Spelling)',
        'format' => 'usfm',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Usfm)
        ->and($entry->markupFormat)->toBe('usfm');
});

it('defaults unknown import formats to sword', function () {
    $entry = CatalogEntry::fromArray([
        'short' => 'KJV',
        'name' => 'King James Version',
        'import_as' => 'invalid',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Sword);
});
