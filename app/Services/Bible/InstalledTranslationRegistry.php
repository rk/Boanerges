<?php

namespace App\Services\Bible;

use App\Data\TranslationConfig;
use Illuminate\Support\Collection;

class InstalledTranslationRegistry
{
    /**
     * @return Collection<int, TranslationConfig>
     */
    public function all(): Collection
    {
        /** @var list<array{id: string, module: string, name: string, abbrev: string}> $translations */
        $translations = config('boanerges.translations');

        return collect($translations)->map(
            fn(array $translation): TranslationConfig => TranslationConfig::fromArray($translation),
        );
    }

    public function find(string $id): TranslationConfig
    {
        $translation = $this->all()->first(
            fn(TranslationConfig $translation): bool => $translation->id === strtolower($id),
        );

        if ($translation === null) {
            abort(404, "Translation \"{$id}\" is not installed.");
        }

        return $translation;
    }
}
