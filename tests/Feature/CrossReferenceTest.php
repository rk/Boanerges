<?php

use App\Jobs\Bible\ImportCrossReferencesJob;
use App\Services\Bible\CrossReferenceService;
use Illuminate\Support\Facades\Bus;
use Illuminate\Support\Facades\Storage;

test('cross reference import job loads bundled data', function (): void {
    $path = Storage::disk('extras')->path('cross-references/cross_references.txt');

    if (! is_file($path)) {
        $this->markTestSkipped('Cross-reference bundle missing.');
    }

    Bus::dispatchSync(new ImportCrossReferencesJob(force: true));

    expect(app(CrossReferenceService::class)->isImported())->toBeTrue();

    $refs = app(CrossReferenceService::class)->forVerse('gen', 1, 1);

    expect($refs)->not->toBeEmpty()
        ->and($refs[0])->toHaveKeys(['rank', 'bookId', 'bookName', 'chapter', 'verse'])
        ->and($refs[0]['bookName'])->toBeString()->not->toBeEmpty();
})->group('crossrefs');

test('cross references are sorted by ascending rank', function (): void {
    $path = Storage::disk('extras')->path('cross-references/cross_references.txt');

    if (! is_file($path)) {
        $this->markTestSkipped('Cross-reference bundle missing.');
    }

    Bus::dispatchSync(new ImportCrossReferencesJob(force: true));

    $refs = app(CrossReferenceService::class)->forVerse('gen', 1, 1);

    expect($refs)->not->toBeEmpty();

    $ranks = array_column($refs, 'rank');

    expect($ranks)->toBe(collect($ranks)->sort()->values()->all());
})->group('crossrefs');

test('cross references resolve common book aliases', function (): void {
    $path = Storage::disk('extras')->path('cross-references/cross_references.txt');

    if (! is_file($path)) {
        $this->markTestSkipped('Cross-reference bundle missing.');
    }

    Bus::dispatchSync(new ImportCrossReferencesJob(force: true));

    $refs = app(CrossReferenceService::class)->forVerse('mark', 1, 1);

    expect($refs)->not->toBeEmpty();
})->group('crossrefs');

test('cross references api returns ranked results', function (): void {
    $path = Storage::disk('extras')->path('cross-references/cross_references.txt');

    if (! is_file($path)) {
        $this->markTestSkipped('Cross-reference bundle missing.');
    }

    Bus::dispatchSync(new ImportCrossReferencesJob(force: true));

    $response = $this->getJson(route('bible.cross-references', [
        'book' => 'gen',
        'chapter' => 1,
        'verse' => 1,
    ]));

    $response->assertSuccessful();
    expect($response->json('references'))->not->toBeEmpty();
})->group('crossrefs');
