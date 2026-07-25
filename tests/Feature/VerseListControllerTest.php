<?php

use App\Services\VerseList\VerseListStore;
use Illuminate\Support\Facades\Storage;

beforeEach(function (): void {
    Storage::fake('local');
});

test('verse lists can be listed when empty', function (): void {
    $this->getJson(route('verse-lists.index'))
        ->assertSuccessful()
        ->assertJsonPath('lists', []);
});

test('verse lists can be created and retrieved', function (): void {
    $response = $this->putJson(route('verse-lists.store'), [
        'title' => 'Favorites',
        'entries' => [
            ['bookId' => 'gen', 'chapter' => 1, 'verse' => 1],
            ['bookId' => 'jhn', 'chapter' => 3, 'verse' => 16],
        ],
    ]);

    $response->assertSuccessful();
    $response->assertJsonPath('list.title', 'Favorites');
    $response->assertJsonCount(2, 'list.entries');

    $id = $response->json('list.id');

    $this->getJson(route('verse-lists.show', ['id' => $id]))
        ->assertSuccessful()
        ->assertJsonPath('list.title', 'Favorites')
        ->assertJsonCount(2, 'list.entries');

    $this->getJson(route('verse-lists.index'))
        ->assertSuccessful()
        ->assertJsonCount(1, 'lists')
        ->assertJsonPath('lists.0.title', 'Favorites');
});

test('verse lists can be updated', function (): void {
    $created = app(VerseListStore::class)->save('Old title', [
        ['bookId' => 'gen', 'chapter' => 1, 'verse' => 1],
    ]);

    $this->putJson(route('verse-lists.update', ['id' => $created['id']]), [
        'title' => 'New title',
        'entries' => [
            ['bookId' => 'mat', 'chapter' => 5, 'verse' => 3],
        ],
    ])
        ->assertSuccessful()
        ->assertJsonPath('list.title', 'New title')
        ->assertJsonPath('list.entries.0.bookId', 'mat');
});

test('verse lists can be deleted', function (): void {
    $created = app(VerseListStore::class)->save('Temporary', [
        ['bookId' => 'gen', 'chapter' => 1, 'verse' => 1],
    ]);

    $this->deleteJson(route('verse-lists.destroy', ['id' => $created['id']]))
        ->assertSuccessful()
        ->assertJsonPath('deleted', true);

    $this->getJson(route('verse-lists.show', ['id' => $created['id']]))
        ->assertNotFound();
});

test('verse list validation rejects invalid entries', function (): void {
    $this->putJson(route('verse-lists.store'), [
        'title' => 'Bad list',
        'entries' => [
            ['bookId' => 'invalid book!', 'chapter' => 0, 'verse' => 0],
        ],
    ])->assertUnprocessable();
});
