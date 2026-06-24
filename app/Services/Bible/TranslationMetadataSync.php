<?php

namespace App\Services\Bible;

use App\Data\CatalogEntry;
use App\Models\Translation;
use App\Services\Bible\Import\SwordConfReader;

class TranslationMetadataSync
{
    public function __construct(
        private SwordConfReader $swordConf,
    ) {}

    public function applyFromCatalog(Translation $translation, CatalogEntry $entry): Translation
    {
        if ($entry->importAs === 'sword') {
            return $this->applyFromSwordConf($translation, $entry);
        }

        $updates = array_filter([
            'about' => $entry->about,
            'source' => $entry->url,
            'format' => $entry->markupFormat,
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
            'format' => $metadata['format'],
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
}
