<?php

use App\Jobs\Bible\ImportWebster1828DictionaryJob;
use App\Services\Dictionary\DictionaryService;
use App\Services\Dictionary\Webster1828Dictionary;
use Illuminate\Support\Facades\Bus;

test('dictionary import job loads fixture data', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    expect(app(DictionaryService::class)->isImported())->toBeTrue();
});

test('dictionary lookup groups variants by part of speech', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    $result = app(Webster1828Dictionary::class)->lookup('grace');

    expect($result['found'])->toBeTrue()
        ->and($result['word'])->toBe('GRACE')
        ->and($result['variants'])->toHaveCount(2)
        ->and($result['variants'][0]['partOfSpeechExpanded'])->toBe('noun')
        ->and($result['variants'][1]['partOfSpeechExpanded'])->toBe('verb');
});

test('dictionary lookup is case insensitive', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    $result = app(Webster1828Dictionary::class)->lookup('GIRD');

    expect($result['found'])->toBeTrue()
        ->and($result['variants'])->not->toBeEmpty();
});

test('dictionary suggest returns distinct prefix matches', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    $graceMatches = app(Webster1828Dictionary::class)->suggest('gr', 10);
    $girdMatches = app(Webster1828Dictionary::class)->suggest('gi', 10);

    expect($graceMatches)->toContain('GRACE')
        ->and($girdMatches)->toContain('GIRD');
});

test('dictionary suggest requires at least two characters', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    expect(app(Webster1828Dictionary::class)->suggest('g', 10))->toBe([]);
});

test('dictionary lookup returns not found for missing words', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    $result = app(Webster1828Dictionary::class)->lookup('zzzznotaword');

    expect($result['found'])->toBeFalse()
        ->and($result['message'])->toContain('not in this dictionary');
});

test('dictionary api returns importing status before import completes', function (): void {
    $response = $this->getJson(route('bible.dictionary.show', ['word' => 'grace']));

    $response->assertStatus(503)
        ->assertJsonPath('status', 'importing');
});

test('dictionary api returns grouped lookup results', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    $response = $this->getJson(route('bible.dictionary.show', ['word' => 'grace']));

    $response->assertSuccessful()
        ->assertJsonPath('found', true)
        ->assertJsonCount(2, 'variants');
});

test('dictionary suggest api returns suggestions', function (): void {
    Bus::dispatchSync(new ImportWebster1828DictionaryJob(force: true));

    $response = $this->getJson(route('bible.dictionary.suggest', ['q' => 'gr']));

    $response->assertSuccessful()
        ->assertJsonPath('suggestions.0', 'GRACE');
});
