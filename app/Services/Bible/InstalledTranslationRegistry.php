<?php

namespace App\Services\Bible;

use App\Data\TranslationConfig;
use Illuminate\Support\Collection;
use Illuminate\Support\Str;

class InstalledTranslationRegistry
{
    public function __construct(
        private BibleModuleManager $modules,
        private TranslationCatalog $catalog,
    ) {}

    /**
     * @return Collection<int, TranslationConfig>
     */
    public function all(): Collection
    {
        $installed = $this->modules->installedModules();

        return collect($installed)->map(function (array $module): TranslationConfig {
            $moduleKey = $module['key'];
            $catalogEntry = $this->catalog->all()->first(
                fn($entry) => strcasecmp($entry->short, $moduleKey) === 0,
            );

            $name = $catalogEntry?->name
                ?? (string) ($module['description'] ?? $moduleKey);

            return new TranslationConfig(
                id: strtolower($moduleKey),
                module: $moduleKey,
                name: $name,
                abbrev: $moduleKey,
                bundled: $module['bundled'],
            );
        })->sortBy('name')->values();
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

    public function findByModule(string $moduleKey): TranslationConfig
    {
        return $this->find(strtolower($moduleKey));
    }

    public function isInstalled(string $moduleKey): bool
    {
        return $this->modules->isModuleInstalled($moduleKey);
    }

    public function isBundled(string $moduleKey): bool
    {
        return in_array(
            Str::upper($moduleKey),
            array_map('strtoupper', config('boanerges.bundled_modules', [])),
            true,
        );
    }
}
