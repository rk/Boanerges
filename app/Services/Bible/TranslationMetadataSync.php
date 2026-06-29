<?php

namespace App\Services\Bible;

use App\Data\CatalogEntry;
use App\Enums\CatalogImportFormat;
use App\Enums\VerseMarkupFormat;
use App\Models\Translation;
use App\Services\Bible\Import\SwordConfReader;

class TranslationMetadataSync
{
    public function __construct(
        private SwordConfReader $swordConf,
    ) {}

    public function applyFromCatalog(Translation $translation, CatalogEntry $entry): Translation
    {
        if ($entry->importAs === CatalogImportFormat::Sword) {
            return $this->applyFromSwordConf($translation, $entry);
        }

        $updates = array_filter([
            'about' => $entry->about,
            'source' => $entry->url,
            'format' => $entry->markupFormat?->value,
        ], fn($value) => $value !== null && $value !== '');

        if ($updates !== []) {
            $translation->update($updates);
        }

        return $translation->fresh();
    }

    public function applyFromSwordConf(Translation $translation, ?CatalogEntry $entry = null): Translation
    {
        $metadata = $this->swordConf->read($translation->abbrev);

        if ($metadata === null) {
            return $translation;
        }

        $updates = array_filter([
            'name' => $metadata['name'],
            'format' => $this->resolveStoredMarkupFormat($metadata['format']),
            'versification' => $metadata['versification'],
            'about' => $metadata['about'],
            'version_string' => $metadata['version_string'],
            'version_date' => $metadata['version_date'],
            'copyright' => $metadata['copyright'],
            'copyright_contact' => $metadata['copyright_contact'],
            'source' => $metadata['source'] ?? $entry?->url,
        ], fn($value) => $value !== null && $value !== '');

        $translation->update($updates);

        return $translation->fresh();
    }

    private function resolveStoredMarkupFormat(?string $raw): ?string
    {
        if ($raw === null || $raw === '') {
            return null;
        }

        $normalized = strtolower($raw);
        $format = VerseMarkupFormat::tryFrom($normalized);

        return $format !== null ? $format->value : $normalized;
    }
}
