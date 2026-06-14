<?php

namespace App\Services;

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
        return $this->settings->get(self::KEY, $this->defaults());
    }

    /**
     * @param  array<string, mixed>  $settings
     * @return array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string}
     */
    public function update(array $settings): array
    {
        /** @var array{activeView: string, bookId: string, chapter: int, translationId: string, translationBId: string} */
        return $this->settings->update(self::KEY, $settings, $this->defaults());
    }
}
