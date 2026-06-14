<?php

namespace App\Services;

use Illuminate\Support\Facades\File;
use Native\Desktop\Facades\Settings;

class AppSettingsRepository
{
    /**
     * @param  array<string, mixed>  $defaults
     * @return array<string, mixed>
     */
    public function get(string $key, array $defaults): array
    {
        $stored = $this->read($key);

        if (! is_array($stored)) {
            return $defaults;
        }

        return array_merge($defaults, $stored);
    }

    /**
     * @param  array<string, mixed>  $settings
     * @param  array<string, mixed>  $defaults
     * @return array<string, mixed>
     */
    public function update(string $key, array $settings, array $defaults): array
    {
        $merged = array_merge($this->get($key, $defaults), $settings);

        $this->write($key, $merged);

        return $merged;
    }

    private function usesNativeStorage(): bool
    {
        return (bool) config('nativephp-internal.running', false);
    }

    private function read(string $key): mixed
    {
        if ($this->usesNativeStorage()) {
            return Settings::get($key);
        }

        return data_get($this->readFile(), $key);
    }

    private function write(string $key, mixed $value): void
    {
        if ($this->usesNativeStorage()) {
            Settings::set($key, $value);

            return;
        }

        $data = $this->readFile();
        data_set($data, $key, $value);
        $this->writeFile($data);
    }

    /**
     * @return array<string, mixed>
     */
    private function readFile(): array
    {
        $path = $this->filePath();

        if (! File::exists($path)) {
            return [];
        }

        $contents = File::get($path);
        $decoded = json_decode($contents, true);

        return is_array($decoded) ? $decoded : [];
    }

    /**
     * @param  array<string, mixed>  $data
     */
    private function writeFile(array $data): void
    {
        File::ensureDirectoryExists(dirname($this->filePath()));

        File::put(
            $this->filePath(),
            json_encode($data, JSON_PRETTY_PRINT | JSON_THROW_ON_ERROR),
        );
    }

    private function filePath(): string
    {
        return storage_path('app/private/settings.json');
    }
}
