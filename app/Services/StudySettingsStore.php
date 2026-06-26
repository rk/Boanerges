<?php

namespace App\Services;

use App\Services\Bible\OsisBookId;

class StudySettingsStore
{
    public const KEY = 'study';

    public function __construct(private AppSettingsRepository $settings) {}

    /**
     * @return array{columnCount: int, columns: list<string>, bookId: string, chapter: int, translationId: string, translationBId: string, translationCId: string}
     */
    public function defaults(): array
    {
        /** @var array{columnCount: int, columns: list<string>, bookId: string, chapter: int, translationId: string, translationBId: string, translationCId: string} $defaults */
        $defaults = config('boanerges.study');

        return $defaults;
    }

    /**
     * @return array{columnCount: int, columns: list<string>, bookId: string, chapter: int, translationId: string, translationBId: string, translationCId: string}
     */
    public function get(): array
    {
        /** @var array<string, mixed> $raw */
        $raw = $this->settings->get(self::KEY, []);
        $raw = is_array($raw) ? $raw : [];

        $settings = array_merge($this->defaults(), $this->migrateFromLegacy($raw));
        $settings['bookId'] = OsisBookId::normalize($settings['bookId']) ?? $settings['bookId'];

        return $settings;
    }

    /**
     * @param  array<string, mixed>  $settings
     * @return array{columnCount: int, columns: list<string>, bookId: string, chapter: int, translationId: string, translationBId: string, translationCId: string}
     */
    public function update(array $settings): array
    {
        if (isset($settings['bookId']) && is_string($settings['bookId'])) {
            $settings['bookId'] = OsisBookId::normalize($settings['bookId']) ?? $settings['bookId'];
        }

        $merged = $this->migrateFromLegacy(array_merge($this->get(), $settings));

        /** @var array{columnCount: int, columns: list<string>, bookId: string, chapter: int, translationId: string, translationBId: string, translationCId: string} */
        return $this->settings->update(self::KEY, $merged, $this->defaults());
    }

    /**
     * @param  array<string, mixed>  $stored
     * @return array{columnCount: int, columns: list<string>, bookId: string, chapter: int, translationId: string, translationBId: string, translationCId: string}
     */
    private function migrateFromLegacy(array $stored): array
    {
        $defaults = $this->defaults();

        if (isset($stored['columnCount'])) {
            $columnCount = (int) $stored['columnCount'];
            $columnCount = in_array($columnCount, [1, 2, 3], true) ? $columnCount : 1;

            $columns = is_array($stored['columns'] ?? null)
                ? $this->sanitizeColumns($columnCount, $stored['columns'])
                : $this->normalizeColumnSlots($columnCount, []);

            return [
                'columnCount' => $columnCount,
                'columns' => $columns,
                'bookId' => (string) ($stored['bookId'] ?? $defaults['bookId']),
                'chapter' => (int) ($stored['chapter'] ?? $defaults['chapter']),
                'translationId' => (string) ($stored['translationId'] ?? $defaults['translationId']),
                'translationBId' => (string) ($stored['translationBId'] ?? $defaults['translationBId']),
                'translationCId' => (string) ($stored['translationCId'] ?? $defaults['translationCId']),
            ];
        }

        $activeView = (string) ($stored['activeView'] ?? 'bible');

        $layout = match ($activeView) {
            'comparison' => ['columnCount' => 2, 'columns' => ['bible-secondary']],
            'scribe' => ['columnCount' => 3, 'columns' => ['scribe', 'bible-secondary']],
            default => ['columnCount' => 1, 'columns' => []],
        };

        return [
            'columnCount' => $layout['columnCount'],
            'columns' => $layout['columns'],
            'bookId' => (string) ($stored['bookId'] ?? $defaults['bookId']),
            'chapter' => (int) ($stored['chapter'] ?? $defaults['chapter']),
            'translationId' => (string) ($stored['translationId'] ?? $defaults['translationId']),
            'translationBId' => (string) ($stored['translationBId'] ?? $defaults['translationBId']),
            'translationCId' => (string) ($stored['translationCId'] ?? $defaults['translationCId']),
        ];
    }

    /**
     * @param  list<mixed>  $columns
     * @return list<string>
     */
    private function sanitizeColumns(int $columnCount, array $columns): array
    {
        $allowed = ['bible-secondary', 'notes', 'scribe', 'search', 'cross-references'];

        $filtered = array_values(array_filter(
            $columns,
            fn($column) => is_string($column) && in_array($column, $allowed, true),
        ));

        return $this->normalizeColumnSlots($columnCount, $filtered);
    }

    /**
     * @param  list<string>  $columns
     * @return list<string>
     */
    private function normalizeColumnSlots(int $columnCount, array $columns): array
    {
        $needed = max(0, $columnCount - 1);
        $normalized = array_slice($columns, 0, $needed);

        while (count($normalized) < $needed) {
            $normalized[] = 'bible-secondary';
        }

        return $normalized;
    }
}
