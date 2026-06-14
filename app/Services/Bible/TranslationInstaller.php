<?php

namespace App\Services\Bible;

use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Storage;
use ZipArchive;

class TranslationInstaller
{
    public function __construct(
        private BibleModuleManager $modules,
        private TranslationCatalog $catalog,
        private InstalledTranslationRegistry $registry,
    ) {}

    /**
     * Install a translation module from the catalog.
     */
    public function install(string $moduleKey): void
    {
        $entry = $this->catalog->find($moduleKey);

        if ($this->registry->isBundled($entry->short)) {
            abort(422, "The {$entry->short} translation is bundled with the app and cannot be installed separately.");
        }

        if ($this->modules->isModuleInstalled($entry->short)) {
            abort(409, "The {$entry->short} translation is already installed.");
        }

        $localRoot = $this->modules->localRoot();
        $zipPath = Storage::disk('local')->path('tmp/' . $entry->short . '.sword-module.zip');

        if (! is_dir(dirname($zipPath))) {
            mkdir(dirname($zipPath), 0755, true);
        }

        $response = Http::timeout(120)->sink($zipPath)->get($entry->url);

        if (! $response->successful() || ! is_file($zipPath) || filesize($zipPath) === 0) {
            @unlink($zipPath);
            abort(502, "Failed to download the {$entry->short} translation.");
        }

        $zip = new ZipArchive();

        if ($zip->open($zipPath) !== true) {
            @unlink($zipPath);
            abort(502, "Failed to open the downloaded archive for {$entry->short}.");
        }

        $zip->extractTo($localRoot);
        $zip->close();

        @unlink($zipPath);

        $this->modules->clearCache();

        if (! $this->modules->isModuleInstalled($entry->short)) {
            abort(502, "Installation failed: {$entry->short} module files were not found after extraction.");
        }

        $this->verify($entry->short);
    }

    public function uninstall(string $moduleKey): void
    {
        $entry = $this->catalog->find($moduleKey);

        if ($this->registry->isBundled($entry->short)) {
            abort(422, "The {$entry->short} translation is bundled with the app and cannot be removed.");
        }

        if (! $this->modules->isModuleInstalled($entry->short)) {
            abort(404, "The {$entry->short} translation is not installed.");
        }

        $parsed = new \rk\PhpSword\SwordModules([$this->modules->localRoot()]);
        $parsed->parseModules();

        $resolvedKey = collect(array_keys($parsed->modules))
            ->first(fn(string $key): bool => strcasecmp($key, $entry->short) === 0);

        if ($resolvedKey === null) {
            abort(422, "The {$entry->short} translation is not installed locally.");
        }

        $moduleFolder = $parsed->modulePaths[$resolvedKey];
        $datapath = (string) ($parsed->modules[$resolvedKey]['datapath'] ?? '');
        $moduleDataPath = $moduleFolder . DIRECTORY_SEPARATOR . $datapath;

        $confPath = $this->findConfFile($this->modules->localRoot(), $entry->short);

        if ($confPath !== null) {
            @unlink($confPath);
        }

        $this->deleteDirectory($moduleDataPath);

        $this->modules->clearCache();
    }

    private function verify(string $moduleKey): void
    {
        $bible = $this->modules->open($moduleKey);

        $text = trim($bible->get(
            books: 'Genesis',
            chapters: 1,
            verses: 1,
            clean: true,
            join: '',
        ));

        if ($text === '') {
            abort(502, "Smoke test failed for {$moduleKey}.");
        }
    }

    private function findConfFile(string $root, string $moduleKey): ?string
    {
        $confDir = $root . DIRECTORY_SEPARATOR . 'mods.d';

        if (! is_dir($confDir)) {
            return null;
        }

        foreach (scandir($confDir) ?: [] as $file) {
            if (strcasecmp(pathinfo($file, PATHINFO_FILENAME), $moduleKey) === 0 && str_ends_with(strtolower($file), '.conf')) {
                return $confDir . DIRECTORY_SEPARATOR . $file;
            }
        }

        return null;
    }

    private function deleteDirectory(string $directory): void
    {
        if (! is_dir($directory)) {
            return;
        }

        foreach (scandir($directory) ?: [] as $item) {
            if ($item === '.' || $item === '..') {
                continue;
            }

            $path = $directory . DIRECTORY_SEPARATOR . $item;

            if (is_dir($path)) {
                $this->deleteDirectory($path);
            } else {
                @unlink($path);
            }
        }

        @rmdir($directory);
    }
}
