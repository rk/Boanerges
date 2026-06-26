<?php

use Illuminate\Support\Facades\Storage;

beforeEach(function (): void {
    Storage::fake('local');
});

test('returns empty content when no notes exist', function (): void {
    $response = $this->getJson(route('notes.chapters.show', ['book' => 'gen', 'chapter' => 1]));

    $response->assertSuccessful();
    $response->assertJsonPath('content', '');
});

test('saves and retrieves chapter notes', function (): void {
    $this->putJson(route('notes.chapters.update', ['book' => 'gen', 'chapter' => 1]), [
        'content' => 'Chapter summary notes.',
    ])
        ->assertSuccessful()
        ->assertJsonPath('content', 'Chapter summary notes.');

    Storage::disk('local')->assertExists('notes/gen/1.txt');

    $this->getJson(route('notes.chapters.show', ['book' => 'gen', 'chapter' => 1]))
        ->assertSuccessful()
        ->assertJsonPath('content', 'Chapter summary notes.');
});

test('rejects invalid book id', function (): void {
    $this->getJson(route('notes.chapters.show', ['book' => 'INVALID', 'chapter' => 1]))
        ->assertUnprocessable();
});
