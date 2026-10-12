<?php

use App\Data\CatalogEntry;
use App\Enums\CatalogImportFormat;
use App\Enums\VerseMarkupFormat;
use Illuminate\Log\Events\MessageLogged;
use Illuminate\Support\Facades\Event;

uses(Tests\TestCase::class);

it('maps OEB catalog entry to USFM import and markup', function () {
    $entry = CatalogEntry::fromArray([
        'short' => 'OEB',
        'name' => 'Open English Bible (US Spelling)',
        'url' => 'https://openenglishbible.org/oeb/2025.6/OEB-2025.6-US.usfm.zip',
        'import_as' => 'usfm',
        'markup_format' => 'usfm',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Usfm)
        ->and($entry->markupFormat)->toBe(VerseMarkupFormat::Usfm);
});

it('defaults USFM verse markup when only import format is given', function () {
    $entry = CatalogEntry::fromArray([
        'short' => 'OEB',
        'name' => 'Open English Bible (US Spelling)',
        'format' => 'usfm',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Usfm)
        ->and($entry->markupFormat)->toBe(VerseMarkupFormat::Usfm);
});

it('defaults unknown import formats to sword', function () {
    Event::fake([MessageLogged::class]);

    $entry = CatalogEntry::fromArray([
        'short' => 'KJV',
        'name' => 'King James Version',
        'import_as' => 'invalid',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Sword);

    Event::assertDispatched(MessageLogged::class, function (MessageLogged $event): bool {
        return $event->level === 'warning'
            && str_contains($event->message, 'import_as');
    });
});

it('logs when legacy format string is unknown', function () {
    Event::fake([MessageLogged::class]);

    $entry = CatalogEntry::fromArray([
        'short' => 'KJV',
        'name' => 'King James Version',
        'format' => 'not-a-real-format',
    ]);

    expect($entry->importAs)->toBe(CatalogImportFormat::Sword);

    Event::assertDispatched(MessageLogged::class, function (MessageLogged $event): bool {
        return $event->level === 'warning'
            && str_contains($event->message, 'import format');
    });
});

it('ignores unknown markup formats from catalog json', function () {
    $entry = CatalogEntry::fromArray([
        'short' => 'KJV',
        'name' => 'King James Version',
        'markup_format' => 'unknown-format',
    ]);

    expect($entry->markupFormat)->toBeNull();
});
