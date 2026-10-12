<?php

use Illuminate\Support\Facades\Storage;

beforeEach(function (): void {
    Storage::fake('local');
});

test('returns empty verses when no draft exists', function (): void {
    $response = $this->getJson(route('scribe.chapters.show', ['book' => 'gen', 'chapter' => 1]));

    $response->assertSuccessful();
    $response->assertJsonPath('verses', []);
});

test('saves and retrieves a chapter draft', function (): void {
    $payload = [
        'verses' => [
            ['verse' => 1, 'text' => 'In the beginning...'],
            ['verse' => 2, 'text' => 'And the earth was...'],
        ],
    ];

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), $payload)
        ->assertSuccessful()
        ->assertJsonPath('verses.0.text', 'In the beginning...');

    Storage::disk('local')->assertExists('scribe/gen/1.json');

    $this->getJson(route('scribe.chapters.show', ['book' => 'gen', 'chapter' => 1]))
        ->assertSuccessful()
        ->assertJsonPath('verses.1.text', 'And the earth was...');
});

test('persists paragraphStart overrides', function (): void {
    $payload = [
        'verses' => [
            ['verse' => 5, 'text' => 'New paragraph here.', 'paragraphStart' => true],
        ],
    ];

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), $payload)
        ->assertSuccessful()
        ->assertJsonPath('verses.0.paragraphStart', true);
});

test('persists text with double line breaks within a verse', function (): void {
    $payload = [
        'verses' => [
            ['verse' => 1, 'text' => "First paragraph.\n\nSecond paragraph."],
        ],
    ];

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), $payload)
        ->assertSuccessful()
        ->assertJsonPath('verses.0.text', "First paragraph.\n\nSecond paragraph.");
});

test('rejects invalid book id', function (): void {
    $this->getJson(route('scribe.chapters.show', ['book' => 'GEN', 'chapter' => 1]))
        ->assertUnprocessable();

    $this->getJson(route('scribe.chapters.show', ['book' => 'toolongbookid', 'chapter' => 1]))
        ->assertUnprocessable();
});

test('rejects invalid payload', function (): void {
    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'verses' => [
            ['text' => 'Missing verse number'],
        ],
    ])->assertUnprocessable();

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'verses' => [
            ['verse' => 1, 'text' => 'Ok', 'paragraphStart' => 'yes'],
        ],
    ])->assertUnprocessable();
});

test('overwrites existing draft on put', function (): void {
    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'verses' => [['verse' => 1, 'text' => 'First draft']],
    ])->assertSuccessful();

    $this->putJson(route('scribe.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'verses' => [['verse' => 1, 'text' => 'Second draft']],
    ])->assertSuccessful()
        ->assertJsonPath('verses.0.text', 'Second draft');
});
