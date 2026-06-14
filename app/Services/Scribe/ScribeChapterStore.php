<?php

namespace App\Services\Scribe;

use Illuminate\Support\Facades\Storage;
use InvalidArgumentException;
use JsonException;

class ScribeChapterStore
{
    /**
     * @return list<array{verse: int, text: string, paragraphStart?: bool}>
     */
    public function get(string $bookId, int $chapter): array
    {
        $this->assertValidBookId($bookId);

        $path = $this->path($bookId, $chapter);

        if (! Storage::disk('local')->exists($path)) {
            return [];
        }

        $contents = Storage::disk('local')->get($path);

        if ($contents === null || $contents === '') {
            return [];
        }

        try {
            $decoded = json_decode($contents, true, 512, JSON_THROW_ON_ERROR);
        } catch (JsonException) {
            return [];
        }

        if (! is_array($decoded)) {
            return [];
        }

        /** @var list<array{verse: int, text: string, paragraphStart?: bool}> $decoded */
        return $decoded;
    }

    /**
     * @param  list<array{verse: int, text: string, paragraphStart?: bool}>  $verses
     */
    public function put(string $bookId, int $chapter, array $verses): void
    {
        $this->assertValidBookId($bookId);

        $path = $this->path($bookId, $chapter);

        Storage::disk('local')->makeDirectory(dirname($path));

        Storage::disk('local')->put(
            $path,
            json_encode($verses, JSON_PRETTY_PRINT | JSON_THROW_ON_ERROR),
        );
    }

    public function path(string $bookId, int $chapter): string
    {
        return "scribe/{$bookId}/{$chapter}.json";
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
