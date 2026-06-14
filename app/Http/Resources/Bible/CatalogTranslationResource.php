<?php

namespace App\Http\Resources\Bible;

use App\Data\CatalogEntry;
use App\Services\Bible\BibleModuleManager;
use App\Services\Bible\InstalledTranslationRegistry;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin CatalogEntry */
class CatalogTranslationResource extends JsonResource
{
    public function __construct(
        CatalogEntry $resource,
        private BibleModuleManager $modules,
        private InstalledTranslationRegistry $registry,
    ) {
        parent::__construct($resource);
    }

    /**
     * @return array{
     *     id: string,
     *     module: string,
     *     name: string,
     *     abbrev: string,
     *     installed: bool,
     *     bundled: bool
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var CatalogEntry $entry */
        $entry = $this->resource;

        return [
            'id' => $entry->id(),
            'module' => $entry->module(),
            'name' => $entry->name,
            'abbrev' => $entry->short,
            'installed' => $this->modules->isModuleInstalled($entry->short),
            'bundled' => $this->registry->isBundled($entry->short),
        ];
    }
}
