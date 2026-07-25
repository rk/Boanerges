<?php

namespace App\Services\Bible\Import;

use App\Enums\CatalogImportFormat;
use App\Models\Translation;
use Closure;

interface TranslationFormatImporter
{
    public function format(): CatalogImportFormat;

    public function isAlreadyAvailable(Translation $translation): bool;

    /**
     * Process a downloaded archive. Returns a path for non-SWORD sources, or null when handled in place.
     */
    public function prepareDownloadedSource(Translation $translation, string $downloadedPath): ?string;

    public function import(Translation $translation, ?string $sourcePath, Closure $onProgress): void;

    public function verify(Translation $translation): void;
}
