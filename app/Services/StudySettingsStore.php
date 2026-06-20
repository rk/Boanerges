<?php

namespace App\Services;

use App\Services\Bible\OsisBookId;

class StudySettingsStore
{
    public const KEY = 'study';

    public function __construct(private AppSettingsRepository $settings) {}

    /**
     * @return array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string}
     */
    public function defaults(): array
    {
        /** @var array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string} $defaults */
        $defaults = config('boanerges.study');

        return $defaults;
    }

    /**
     * @return array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string}
     */
    public function get(): array
    {
        /** @var array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string} */
        $settings = $this->settings->get(self::KEY, $this->defaults());
        $settings['bookId'] = OsisBookId::normalize($settings['bookId']) ?? $settings['bookId'];

        return $settings;
    }

    /**
     * @param  array<string, mixed>  $settings
     * @return array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string}
     */
    public function update(array $settings): array
    {
        if (isset($settings['bookId']) && is_string($settings['bookId'])) {
            $settings['bookId'] = OsisBookId::normalize($settings['bookId']) ?? $settings['bookId'];
        }

        /** @var array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string} */
        return $this->settings->update(self::KEY, $settings, $this->defaults());
    }
}
