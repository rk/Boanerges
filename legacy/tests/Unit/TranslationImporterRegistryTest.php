<?php

use App\Enums\CatalogImportFormat;
use App\Services\Bible\Import\TranslationImporterRegistry;

uses(Tests\TestCase::class);

test('every catalog import format resolves to an importer', function () {
    $registry = app(TranslationImporterRegistry::class);

    foreach (CatalogImportFormat::cases() as $format) {
        expect($registry->get($format)->format())->toBe($format);
    }
});
