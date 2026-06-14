<?php

namespace App\Services\Bible;

use App\Data\CatalogEntry;
use Illuminate\Support\Collection;
use Illuminate\Support\Facades\Storage;

class TranslationCatalog
{
    /**
     * @return Collection<int, CatalogEntry>
     */
    public function all(): Collection
    {
        $path = config('boanerges.catalog_path');
        $contents = Storage::disk('extras')->get($path);

        if ($contents === null) {
            return collect();
        }

        /** @var list<array{short: string, name: string, url: string}> $entries */
        $entries = json_decode($contents, true, flags: JSON_THROW_ON_ERROR);

        return collect($entries)->map(
            fn(array $entry): CatalogEntry => CatalogEntry::fromArray($entry),
        );
    }

    public function find(string $moduleKey): CatalogEntry
    {
        $entry = $this->all()->first(
            fn(CatalogEntry $entry): bool => strcasecmp($entry->short, $moduleKey) === 0,
        );

        if ($entry === null) {
            abort(404, "Translation \"{$moduleKey}\" is not in the catalog.");
        }

        return $entry;
    }

    public function findById(string $id): CatalogEntry
    {
        $entry = $this->all()->first(
            fn(CatalogEntry $entry): bool => $entry->id() === strtolower($id),
        );

        if ($entry === null) {
            abort(404, "Translation \"{$id}\" is not in the catalog.");
        }

        return $entry;
    }
}
