<?php

use App\Services\Study\StudyPrintHtmlBuilder;
use Illuminate\Support\Facades\Storage;
use Native\Desktop\Facades\System;

beforeEach(function (): void {
    seedBundledAsvForTests();
    Storage::fake('local');
});

function studyPrintPayload(array $overrides = []): array
{
    return array_merge([
        'includeUserWork' => false,
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

    $this->postJson(route('study.print'), studyPrintPayload())
        ->assertStatus(503);
});

test('prints portrait layout for a single column', function (): void {
    config(['nativephp-internal.running' => true]);

    System::shouldReceive('print')
        ->once()
        ->withArgs(function (string $html, $printer, array $settings): bool {
            expect($html)->toContain('Genesis 1 (ASV)')
                ->and($html)->toContain('beginning')
                ->and($settings)->toMatchArray([
                    'pageSize' => 'A4',
                    'landscape' => false,
                ]);

            return true;
        });

    $this->postJson(route('study.print'), studyPrintPayload())
        ->assertNoContent();
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
                ->and($html)->toContain('class="line"')
                ->and($html)->not->toContain('My private scribe draft')
                ->and($settings['landscape'])->toBeTrue();

            return true;
        });

    $this->postJson(route('study.print'), studyPrintPayload([
        'columnCount' => 2,
        'columns' => ['scribe'],
    ]))->assertNoContent();
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
        ->and($blankNotes)->toContain('lined-block')
        ->and($blankNotes)->not->toContain('Chapter notes here.');
});
