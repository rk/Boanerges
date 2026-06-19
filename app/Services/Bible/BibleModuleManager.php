<?php

namespace App\Services\Bible;

use App\Exceptions\BibleModuleNotInstalledException;
use Illuminate\Support\Facades\Storage;
use rk\PhpSword\SwordBible;
use rk\PhpSword\SwordModules;

class BibleModuleManager
{
    /** @var array<string, SwordBible> */
    private array $openModules = [];

    private ?SwordModules $swordModules = null;

    /**
     * @return list<string>
     */
    public function swordRoots(): array
    {
        $roots = [
            Storage::disk('extras')->path(config('boanerges.bundled_sword_path')),
        ];

        $localRoot = Storage::disk('local')->path(config('boanerges.local_sword_path'));

        if (! is_dir($localRoot)) {
            Storage::disk('local')->makeDirectory(config('boanerges.local_sword_path'));
        }

        $roots[] = $localRoot;

        return $roots;
    }

    public function bundledRoot(): string
    {
        return Storage::disk('extras')->path(config('boanerges.bundled_sword_path'));
    }

    public function localRoot(): string
    {
        return Storage::disk('local')->path(config('boanerges.local_sword_path'));
    }

    public function isInstalled(): bool
    {
        foreach ($this->swordRoots() as $root) {
            if (is_dir($root . DIRECTORY_SEPARATOR . 'mods.d')) {
                return true;
            }
        }

        return false;
    }

    public function isModuleInstalled(string $moduleKey): bool
    {
        return $this->resolveModuleKey($moduleKey) !== null;
    }

    public function resolveModuleKey(string $moduleKey): ?string
    {
        $modules = $this->modules();

        foreach (array_keys($modules->modules) as $key) {
            if (strcasecmp($key, $moduleKey) === 0) {
                return $key;
            }
        }

        return null;
    }

    /**
     * @return list<array{key: string, description: string|null, bundled: bool}>
     */
    public function installedModules(): array
    {
        $parsed = $this->modules();
        $bundledRoot = realpath($this->bundledRoot()) ?: $this->bundledRoot();
        $installed = [];

        foreach ($parsed->modules as $moduleKey => $metadata) {
            $modulePath = $parsed->modulePaths[$moduleKey] ?? '';
            $resolvedPath = realpath($modulePath) ?: $modulePath;
            $bundled = str_starts_with($resolvedPath, $bundledRoot);

            $installed[] = [
                'key' => $moduleKey,
                'description' => isset($metadata['description']) ? (string) $metadata['description'] : null,
                'bundled' => $bundled,
            ];
        }

        return $installed;
    }

    public function open(string $moduleKey): SwordBible
    {
        $resolvedKey = $this->resolveModuleKey($moduleKey) ?? $moduleKey;

        if (isset($this->openModules[$resolvedKey])) {
            return $this->openModules[$resolvedKey];
        }

        if (! $this->isInstalled()) {
            throw BibleModuleNotInstalledException::missing($moduleKey);
        }

        $modules = $this->modules();

        if ($this->resolveModuleKey($moduleKey) === null) {
            throw BibleModuleNotInstalledException::missing($moduleKey);
        }

        $this->openModules[$resolvedKey] = $modules->getBibleFromModule($resolvedKey);

        return $this->openModules[$resolvedKey];
    }

    public function clearCache(): void
    {
        $this->openModules = [];
        $this->swordModules = null;
    }

    private function modules(): SwordModules
    {
        if ($this->swordModules === null) {
            $this->swordModules = new SwordModules($this->swordRoots());
            $this->swordModules->parseModules();
        }

        return $this->swordModules;
    }
}
