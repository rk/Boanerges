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
        $modules = $this->modules();

        return array_key_exists($moduleKey, $modules->modules);
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
        if (isset($this->openModules[$moduleKey])) {
            return $this->openModules[$moduleKey];
        }

        if (! $this->isInstalled()) {
            throw BibleModuleNotInstalledException::missing($moduleKey);
        }

        $modules = $this->modules();

        if (! array_key_exists($moduleKey, $modules->modules)) {
            throw BibleModuleNotInstalledException::missing($moduleKey);
        }

        $this->openModules[$moduleKey] = $modules->getBibleFromModule($moduleKey);

        return $this->openModules[$moduleKey];
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
