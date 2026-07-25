<?php

namespace App\Services\VerseList;

use App\Services\Bible\OsisBookId;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;
use InvalidArgumentException;

class VerseListStore
{
    private const INDEX_PATH = 'verse-lists/index.json';

    /**
     * @return list<array{id: string, title: string, updatedAt: string}>
     */
    public function list(): array
    {
        $index = $this->readIndex();

        usort($index, fn(array $a, array $b): int => strcmp($b['updatedAt'], $a['updatedAt']));

        return $index;
    }

    /**
     * @return array{id: string, title: string, entries: list<array{bookId: string, chapter: int, verse: int}>, updatedAt: string}
     */
    public function get(string $id): array
    {
        $path = $this->path($id);

        if (! Storage::disk('local')->exists($path)) {
            throw new InvalidArgumentException("Verse list not found: {$id}");
        }

        /** @var array{id?: string, title?: string, entries?: list<array{bookId?: string, chapter?: int, verse?: int}>, updatedAt?: string} $data */
        $data = json_decode(Storage::disk('local')->get($path) ?? '', true) ?? [];

        return $this->normalizeList($id, $data);
    }

    /**
     * @param  list<array{bookId: string, chapter: int, verse: int}>  $entries
     * @return array{id: string, title: string, entries: list<array{bookId: string, chapter: int, verse: int}>, updatedAt: string}
     */
    public function save(string $title, array $entries, ?string $id = null): array
    {
        $id = $id !== null && $id !== '' ? $id : (string) Str::uuid();
        $normalizedEntries = $this->normalizeEntries($entries);
        $updatedAt = now()->toIso8601String();

        $list = [
            'id' => $id,
            'title' => trim($title) !== '' ? trim($title) : 'Untitled list',
            'entries' => $normalizedEntries,
            'updatedAt' => $updatedAt,
        ];

        Storage::disk('local')->makeDirectory('verse-lists');
        Storage::disk('local')->put($this->path($id), json_encode($list, JSON_THROW_ON_ERROR | JSON_PRETTY_PRINT));

        $index = $this->readIndex();
        $index = array_values(array_filter($index, fn(array $item): bool => $item['id'] !== $id));
        $index[] = [
            'id' => $id,
            'title' => $list['title'],
            'updatedAt' => $updatedAt,
        ];

        $this->writeIndex($index);

        return $list;
    }

    public function delete(string $id): void
    {
        Storage::disk('local')->delete($this->path($id));

        $index = array_values(array_filter(
            $this->readIndex(),
            fn(array $item): bool => $item['id'] !== $id,
        ));

        $this->writeIndex($index);
    }

    private function path(string $id): string
    {
        if (! preg_match('/^[a-zA-Z0-9-]{1,36}$/', $id)) {
            throw new InvalidArgumentException("Invalid verse list id: {$id}");
        }

        return "verse-lists/{$id}.json";
    }

    /**
     * @param  list<array{bookId?: string, chapter?: int, verse?: int}>  $entries
     * @return list<array{bookId: string, chapter: int, verse: int}>
     */
    private function normalizeEntries(array $entries): array
    {
        $normalized = [];

        foreach ($entries as $entry) {
            $bookId = OsisBookId::normalize((string) ($entry['bookId'] ?? ''));

            if ($bookId === null || ! preg_match('/^[a-z0-9]{1,10}$/', $bookId)) {
                throw new InvalidArgumentException('Invalid book id in verse list entry.');
            }

            $chapter = (int) ($entry['chapter'] ?? 0);
            $verse = (int) ($entry['verse'] ?? 0);

            if ($chapter < 1 || $verse < 1) {
                throw new InvalidArgumentException('Invalid chapter or verse in verse list entry.');
            }

            $normalized[] = [
                'bookId' => $bookId,
                'chapter' => $chapter,
                'verse' => $verse,
            ];
        }

        return $normalized;
    }

    /**
     * @param  array{id?: string, title?: string, entries?: list<array{bookId?: string, chapter?: int, verse?: int}>, updatedAt?: string}  $data
     * @return array{id: string, title: string, entries: list<array{bookId: string, chapter: int, verse: int}>, updatedAt: string}
     */
    private function normalizeList(string $id, array $data): array
    {
        return [
            'id' => $id,
            'title' => (string) ($data['title'] ?? 'Untitled list'),
            'entries' => $this->normalizeEntries($data['entries'] ?? []),
            'updatedAt' => (string) ($data['updatedAt'] ?? now()->toIso8601String()),
        ];
    }

    /**
     * @return list<array{id: string, title: string, updatedAt: string}>
     */
    private function readIndex(): array
    {
        if (! Storage::disk('local')->exists(self::INDEX_PATH)) {
            return [];
        }

        /** @var list<array{id?: string, title?: string, updatedAt?: string}> $index */
        $index = json_decode(Storage::disk('local')->get(self::INDEX_PATH) ?? '', true) ?? [];

        return array_values(array_filter(array_map(function (array $item): ?array {
            if (! isset($item['id'], $item['title'], $item['updatedAt'])) {
                return null;
            }

            return [
                'id' => (string) $item['id'],
                'title' => (string) $item['title'],
                'updatedAt' => (string) $item['updatedAt'],
            ];
        }, $index)));
    }

    /**
     * @param  list<array{id: string, title: string, updatedAt: string}>  $index
     */
    private function writeIndex(array $index): void
    {
        Storage::disk('local')->makeDirectory('verse-lists');
        Storage::disk('local')->put(self::INDEX_PATH, json_encode($index, JSON_THROW_ON_ERROR | JSON_PRETTY_PRINT));
    }
}
