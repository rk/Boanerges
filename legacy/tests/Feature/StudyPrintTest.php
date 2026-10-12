<?php

use App\Services\Study\StudyPrintHtmlBuilder;
use App\Services\Study\StudyPrintService;
use Illuminate\Support\Facades\Storage;
use Native\Desktop\DataObjects\Printer;
use Native\Desktop\Facades\System;

beforeEach(function (): void {
    seedBundledAsvForTests();
    Storage::fake('local');
});

function studyPrintPayload(array $overrides = []): array
{
    return array_merge([
        'includeUserWork' => false,
        'printerName' => null,
        'columnCount' => 1,
        'columns' => [],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ], $overrides);
}

test('rejects invalid print payload', function (): void {
    $this->postJson(route('study.print'), [
        'includeUserWork' => 'yes',
        'columnCount' => 2,
        'columns' => [],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ])->assertUnprocessable();
});

test('returns service unavailable outside the desktop app', function (): void {
    config(['nativephp-internal.running' => false]);

    $this->getJson(route('study.printers.index'))
        ->assertStatus(503);

    $this->postJson(route('study.print'), studyPrintPayload())
        ->assertStatus(503);
});

test('lists available printers in the desktop app', function (): void {
    config(['nativephp-internal.running' => true]);

    System::shouldReceive('printers')
        ->once()
        ->andReturn([
            new Printer('Brother_QL', 'Brother QL', 'Local printer', []),
        ]);

    $this->getJson(route('study.printers.index'))
        ->assertSuccessful()
        ->assertJsonPath('printers.0.name', 'Brother_QL')
        ->assertJsonPath('printers.0.displayName', 'Brother QL');
});

test('prints portrait layout for a single column', function (): void {
    config(['nativephp-internal.running' => true]);

    System::shouldReceive('print')
        ->once()
        ->withArgs(function (string $html, $printer, array $settings): bool {
            expect($html)->toContain('Genesis 1 (ASV)')
                ->and($html)->toContain('beginning')
                ->and($printer)->toBeNull()
                ->and($settings)->toMatchArray([
                    'landscape' => false,
                    'pageSize' => 'A4',
                    'silent' => false,
                    'usePrinterDefaultPageSize' => false,
                    'printBackground' => true,
                ]);

            return true;
        });

    $this->postJson(route('study.print'), studyPrintPayload())
        ->assertNoContent();
});

test('prints to a selected printer with explicit page size and orientation', function (): void {
    config(['nativephp-internal.running' => true]);

    $printer = new Printer('Brother_QL', 'Brother QL', 'Local printer', []);

    System::shouldReceive('printers')
        ->once()
        ->andReturn([$printer]);

    System::shouldReceive('print')
        ->once()
        ->withArgs(function (string $html, ?Printer $selectedPrinter, array $settings) use ($printer): bool {
            expect($selectedPrinter?->name)->toBe($printer->name)
                ->and($settings)->toMatchArray([
                    'landscape' => false,
                    'pageSize' => 'A4',
                    'silent' => false,
                    'usePrinterDefaultPageSize' => false,
                    'printBackground' => true,
                ]);

            return true;
        });

    $this->postJson(route('study.print'), studyPrintPayload([
        'printerName' => 'Brother_QL',
    ]))->assertNoContent();
});

test('rejects an unavailable printer', function (): void {
    config(['nativephp-internal.running' => true]);

    System::shouldReceive('printers')
        ->once()
        ->andReturn([]);

    $this->postJson(route('study.print'), studyPrintPayload([
        'printerName' => 'Missing_Printer',
    ]))->assertUnprocessable();
});

test('prints landscape layout with scribe lined area and no user draft text', function (): void {
    config(['nativephp-internal.running' => true]);

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'verses' => [['verse' => 1, 'text' => 'My private scribe draft']],
    ])->assertSuccessful();

    System::shouldReceive('print')
        ->once()
        ->withArgs(function (string $html, $printer, array $settings): bool {
            expect($html)->toContain('Scribe')
                ->and($html)->toContain('class="lined-block"')
                ->and($html)->toContain('data:image/svg+xml')
                ->and($html)->toContain('background-size: 100% 30.6px')
                ->and($html)->not->toContain('scribe-ghost')
                ->and($html)->not->toMatch('/<div class="column column-scribe">[\s\S]*?<\/div>\s*<\/div>\s*<\/div>.*beginning/s')
                ->and($html)->not->toContain('My private scribe draft')
                ->and($settings['landscape'])->toBeTrue()
                ->and($settings['pageSize'])->toBe('A4')
                ->and($settings['usePrinterDefaultPageSize'])->toBeFalse()
                ->and($settings['printBackground'])->toBeTrue();

            return true;
        });

    $this->postJson(route('study.print'), studyPrintPayload([
        'columnCount' => 2,
        'columns' => ['scribe'],
    ]))->assertNoContent();
});

test('exports html when html destination selected', function (): void {
    config(['nativephp-internal.running' => true]);

    $this->mock(StudyPrintService::class, function ($mock): void {
        $mock->shouldReceive('print')
            ->once()
            ->withArgs(function (array $study, bool $includeUserWork, ?string $printerName): bool {
                expect($printerName)->toBe(StudyPrintService::HTML_DESTINATION)
                    ->and($includeUserWork)->toBeFalse();

                return true;
            })
            ->andReturn('/tmp/Genesis 1 study.html');
    });

    $this->postJson(route('study.print'), studyPrintPayload([
        'printerName' => StudyPrintService::HTML_DESTINATION,
    ]))
        ->assertSuccessful()
        ->assertJsonPath('path', '/tmp/Genesis 1 study.html');
});

test('includes notes when requested and otherwise prints lined notes', function (): void {
    $builder = app(StudyPrintHtmlBuilder::class);

    $this->putJson(route('notes.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'content' => 'Chapter notes here.',
    ])->assertSuccessful();

    $withNotes = $builder->build([
        'columnCount' => 2,
        'columns' => ['notes'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ], true);

    $blankNotes = $builder->build([
        'columnCount' => 2,
        'columns' => ['notes'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ], false);

    expect($withNotes)->toContain('Chapter notes here.')
        ->and($blankNotes)->toContain('class="lined-block"')
        ->and($blankNotes)->toContain('data:image/svg+xml')
        ->and($blankNotes)->not->toContain('Chapter notes here.');
});

test('includes scribe content when requested and otherwise prints lined scribe', function (): void {
    $builder = app(StudyPrintHtmlBuilder::class);

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'verses' => [['verse' => 1, 'text' => 'My scribe reading draft']],
    ])->assertSuccessful();

    $withScribe = $builder->build([
        'columnCount' => 2,
        'columns' => ['scribe'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ], true);

    $blankScribe = $builder->build([
        'columnCount' => 2,
        'columns' => ['scribe'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ], false);

    expect($withScribe)->toContain('My scribe reading draft')
        ->and($withScribe)->toContain('column-scribe-content')
        ->and($withScribe)->not->toContain('class="lined-block"')
        ->and($blankScribe)->toContain('class="lined-block"')
        ->and($blankScribe)->not->toContain('My scribe reading draft');
});

test('prints verse list entries with verse text', function (): void {
    $builder = app(StudyPrintHtmlBuilder::class);

    $html = $builder->build([
        'columnCount' => 2,
        'columns' => ['verse-list'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
        'verseList' => [
            'title' => 'My Passages',
            'entries' => [
                ['bookId' => 'gen', 'chapter' => 1, 'verse' => 1],
            ],
        ],
    ], false);

    expect($html)->toContain('My Passages')
        ->and($html)->toContain('Genesis 1:1')
        ->and($html)->toContain('beginning')
        ->and($html)->toContain('column-verse-list');
});

test('comparison column prints interactive placeholder', function (): void {
    $builder = app(StudyPrintHtmlBuilder::class);

    $html = $builder->build([
        'columnCount' => 2,
        'columns' => ['comparison'],
        'bookId' => 'gen',
        'chapter' => 1,
        'translationId' => 'asv',
        'translationBId' => 'asv',
        'translationCId' => 'asv',
    ], false);

    expect($html)->toContain('Comparison')
        ->and($html)->toContain('Interactive view — not included in print.');
});
