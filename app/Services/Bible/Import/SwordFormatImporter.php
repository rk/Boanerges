<?php

namespace App\Services\Bible\Import;

use App\Enums\CatalogImportFormat;
use App\Models\Translation;
use App\Services\Bible\BibleModuleManager;
use Closure;
use ZipArchive;

class SwordFormatImporter implements TranslationFormatImporter
{
    public function __construct(
        private BibleModuleManager $modules,
        private SwordImporter $swordImporter,
    ) {}

    public function format(): CatalogImportFormat
    {
        return CatalogImportFormat::Sword;
    }

    public function isAlreadyAvailable(Translation $translation): bool
    {
        return $this->modules->isModuleInstalled($translation->abbrev);
    }

    public function prepareDownloadedSource(Translation $translation, string $downloadedPath): ?string
    {
        $zip = new ZipArchive();

        if ($zip->open($downloadedPath) !== true) {
            @unlink($downloadedPath);
            throw new \RuntimeException("Failed to open archive for {$translation->abbrev}.");
        }

        $zip->extractTo($this->modules->localRoot());
        $zip->close();
        @unlink($downloadedPath);
        $this->modules->clearCache();

        return null;
    }

    public function import(Translation $translation, ?string $sourcePath, Closure $onProgress): void
    {
        $this->swordImporter->progressConfigure(20, 50, $onProgress);
        $this->swordImporter->import($translation->abbrev, $translation->abbrev);
    }

    public function verify(Translation $translation): void
    {
        $this->swordImporter->verify($translation->abbrev);
    }
}
