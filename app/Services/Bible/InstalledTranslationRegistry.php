<?php

namespace App\Services\Bible;

use App\Data\TranslationConfig;
use App\Enums\TranslationInstallStatus;
use App\Models\Translation;
use Illuminate\Support\Collection;
use Illuminate\Support\Str;

class InstalledTranslationRegistry
{
    /**
     * @return Collection<int, TranslationConfig>
     */
    public function all(): Collection
    {
        return Translation::query()
            ->where('install_status', TranslationInstallStatus::Ready)
            ->orderBy('name')
            ->get()
            ->map(fn(Translation $translation): TranslationConfig => $this->mapConfig($translation))
            ->values();
    }

    public function find(string $id): TranslationConfig
    {
        $translation = Translation::query()
            ->where('abbrev', strtolower($id))
            ->where('install_status', TranslationInstallStatus::Ready)
            ->first();

        if ($translation === null) {
            abort(404, "Translation \"{$id}\" is not installed.");
        }

        return $this->mapConfig($translation);
    }

    public function findByModule(string $moduleKey): TranslationConfig
    {
        return $this->find(strtolower($moduleKey));
    }

    public function isInstalled(string $moduleKey): bool
    {
        return Translation::query()
            ->where('abbrev', strtolower($moduleKey))
            ->where('install_status', TranslationInstallStatus::Ready)
            ->exists();
    }

    public function isBundled(string $moduleKey): bool
    {
        return in_array(
            Str::upper($moduleKey),
            array_map('strtoupper', config('boanerges.bundled_modules', [])),
            true,
        );
    }

    public function findModel(string $abbrev): ?Translation
    {
        return Translation::query()->where('abbrev', strtolower($abbrev))->first();
    }

    public function toConfig(Translation $translation): TranslationConfig
    {
        return new TranslationConfig(
            id: $translation->abbrev,
            module: strtoupper($translation->abbrev),
            name: $translation->name,
            abbrev: strtoupper($translation->abbrev),
            bundled: $translation->bundled,
            about: $translation->about,
            installStatus: $translation->installStatusValue(),
            installStep: $translation->install_step,
            installError: $translation->install_error,
        );
    }

    private function mapConfig(Translation $translation): TranslationConfig
    {
        return $this->toConfig($translation);
    }
}
