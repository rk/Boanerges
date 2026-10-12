<?php

namespace App\Services\Bible\Import;

use App\Services\Bible\BibleModuleManager;

class SwordConfReader
{
    public function __construct(
        private BibleModuleManager $modules,
    ) {}

    /**
     * @return array{
     *     name: string|null,
     *     format: string|null,
     *     versification: string|null,
     *     about: string|null,
     *     version_string: string|null,
     *     version_date: string|null,
     *     copyright: string|null,
     *     copyright_contact: string|null,
     *     source: string|null
     * }|null
     */
    public function read(string $moduleKey): ?array
    {
        $path = $this->confPath($moduleKey);

        if ($path === null) {
            return null;
        }

        return $this->parseFile($path);
    }

    /**
     * @return array{
     *     name: string|null,
     *     format: string|null,
     *     versification: string|null,
     *     about: string|null,
     *     version_string: string|null,
     *     version_date: string|null,
     *     copyright: string|null,
     *     copyright_contact: string|null,
     *     source: string|null
     * }
     */
    public function parseFile(string $path): array
    {
        $raw = $this->parseRaw($path);

        $about = $raw['About'] ?? null;

        if ($about !== null) {
            $about = trim(str_replace(['\\par', '\par'], "\n\n", $about));
        }

        return [
            'name' => $raw['Description'] ?? null,
            'format' => isset($raw['SourceType']) ? strtolower($raw['SourceType']) : null,
            'versification' => $raw['Versification'] ?? null,
            'about' => $about,
            'version_string' => $raw['Version'] ?? null,
            'version_date' => $raw['SwordVersionDate'] ?? $raw['RevDate'] ?? null,
            'copyright' => $raw['DistributionLicense'] ?? $raw['Copyright'] ?? null,
            'copyright_contact' => $raw['CopyrightHolder'] ?? $raw['CopyrightContact'] ?? null,
            'source' => $raw['TextSource'] ?? null,
        ];
    }

    /** @return array<string, string> */
    private function parseRaw(string $path): array
    {
        $data = [];

        foreach (file($path, FILE_IGNORE_NEW_LINES) ?: [] as $line) {
            $line = trim($line);

            if ($line === '' || str_starts_with($line, '#') || str_starts_with($line, ';')) {
                continue;
            }

            if (str_starts_with($line, '[')) {
                continue;
            }

            if (! str_contains($line, '=')) {
                continue;
            }

            [$key, $value] = explode('=', $line, 2);
            $data[trim($key)] = trim($value);
        }

        return $data;
    }

    private function confPath(string $moduleKey): ?string
    {
        $resolved = $this->modules->resolveModuleKey($moduleKey);

        if ($resolved === null) {
            return null;
        }

        foreach ($this->modules->swordRoots() as $root) {
            $confDir = $root . DIRECTORY_SEPARATOR . 'mods.d';

            if (! is_dir($confDir)) {
                continue;
            }

            foreach (scandir($confDir) ?: [] as $file) {
                if (strcasecmp(pathinfo($file, PATHINFO_FILENAME), $resolved) !== 0) {
                    continue;
                }

                if (! str_ends_with(strtolower($file), '.conf')) {
                    continue;
                }

                return $confDir . DIRECTORY_SEPARATOR . $file;
            }
        }

        return null;
    }
}
