<?php

namespace App\Services\Notes;

use Illuminate\Support\Facades\Storage;
use InvalidArgumentException;

class NotesChapterStore
{
    public function get(string $bookId, int $chapter): string
    {
        $this->assertValidBookId($bookId);

        $path = $this->path($bookId, $chapter);

        if (! Storage::disk('local')->exists($path)) {
            return '';
        }

        return Storage::disk('local')->get($path) ?? '';
    }

    public function put(string $bookId, int $chapter, string $content): void
    {
        $this->assertValidBookId($bookId);

        $path = $this->path($bookId, $chapter);

        Storage::disk('local')->makeDirectory(dirname($path));
        Storage::disk('local')->put($path, $content);
    }

    public function path(string $bookId, int $chapter): string
    {
        return "notes/{$bookId}/{$chapter}.txt";
    }

    public function isValidBookId(string $bookId): bool
    {
        return (bool) preg_match('/^[a-z0-9]{1,10}$/', $bookId);
    }

    private function assertValidBookId(string $bookId): void
    {
        if (! $this->isValidBookId($bookId)) {
            throw new InvalidArgumentException("Invalid book id: {$bookId}");
        }
    }
}
