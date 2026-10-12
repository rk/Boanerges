<?php

namespace App\Services;

class ReadabilitySettingsStore
{
    public const KEY = 'readability';

    public function __construct(private AppSettingsRepository $settings) {}

    /**
     * @return array{fontSize: int, lineHeight: float, theme: string, fontFamily: string, justifyText: bool}
     */
    public function defaults(): array
    {
        /** @var array{fontSize: int, lineHeight: float, theme: string, fontFamily: string, justifyText: bool} $defaults */
        $defaults = config('boanerges.readability');

        return $defaults;
    }

    /**
     * @return array{fontSize: int, lineHeight: float, theme: string, fontFamily: string, justifyText: bool}
     */
    public function get(): array
    {
        /** @var array{fontSize: int, lineHeight: float, theme: string, fontFamily: string, justifyText: bool} $settings */
        $settings = $this->settings->get(self::KEY, $this->defaults());

        if ($settings['theme'] === 'sepia') {
            $settings['theme'] = 'caramellatte';
        }

        return $settings;
    }

    /**
     * @param  array<string, mixed>  $settings
     * @return array{fontSize: int, lineHeight: float, theme: string, fontFamily: string, justifyText: bool}
     */
    public function update(array $settings): array
    {
        /** @var array{fontSize: int, lineHeight: float, theme: string, fontFamily: string, justifyText: bool} */
        return $this->settings->update(self::KEY, $settings, $this->defaults());
    }
}
