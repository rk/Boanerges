<?php

namespace App\Services\Bible\Import;

use App\Enums\CatalogImportFormat;
use InvalidArgumentException;

class TranslationImporterRegistry
{
    /** @var array<string, TranslationFormatImporter> */
    private array $importers = [];

    public function register(TranslationFormatImporter $importer): void
    {
        $this->importers[$importer->format()->value] = $importer;
    }

    public function get(CatalogImportFormat $format): TranslationFormatImporter
    {
        $importer = $this->importers[$format->value] ?? null;

        if ($importer === null) {
            throw new InvalidArgumentException("No importer registered for format: {$format->value}");
        }

        return $importer;
    }

    /**
     * @return list<TranslationFormatImporter>
     */
    public function all(): array
    {
        return array_values($this->importers);
    }
}
