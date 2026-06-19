<?php

namespace App\Http\Resources\Bible;

use App\Data\CatalogEntry;
use App\Services\Bible\InstalledTranslationRegistry;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/** @mixin CatalogEntry */
class CatalogTranslationResource extends JsonResource
{
    public function __construct(
        CatalogEntry $resource,
        private readonly InstalledTranslationRegistry $registry,
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
     *     bundled: bool,
     *     about: string|null,
     *     install_status: string|null
     * }
     */
    public function toArray(Request $request): array
    {
        /** @var CatalogEntry $entry */
        $entry = $this->resource;
        $model = $this->registry->findModel($entry->short);

        return [
            'id' => $entry->id(),
            'module' => $entry->module(),
            'name' => $entry->name,
            'abbrev' => $entry->short,
            'installed' => $model?->isReady() ?? false,
            'bundled' => $this->registry->isBundled($entry->short),
            'about' => $entry->about,
            'install_status' => $model?->install_status->value,
        ];
    }
}
