<?php

namespace App\Services\Bible;

use App\Data\CatalogEntry;
use App\Enums\TranslationInstallStatus;
use App\Enums\TranslationInstallStep;
use App\Jobs\Bible\InstallTranslationJob;
use App\Models\Translation;
use Illuminate\Support\Facades\Bus;

class TranslationInstaller
{
    public function __construct(
        private TranslationCatalog $catalog,
        private InstalledTranslationRegistry $registry,
    ) {}

    public function install(string $moduleKey): Translation
    {
        $entry = $this->catalog->find($moduleKey);

        if ($this->registry->isBundled($entry->short) && $this->registry->isInstalled($entry->short)) {
            abort(422, "The {$entry->short} translation is bundled with the app and cannot be installed separately.");
        }

        $existing = Translation::query()->where('abbrev', strtolower($entry->short))->first();

        if ($existing?->isReady()) {
            abort(409, "The {$entry->short} translation is already installed.");
        }

        if ($existing !== null && ! in_array($existing->install_status, [TranslationInstallStatus::Failed, TranslationInstallStatus::Pending], true)) {
            abort(409, "The {$entry->short} translation is already being installed.");
        }

        $translation = $existing ?? $this->createTranslationRecord($entry);

        if ($existing !== null) {
            $translation->update($this->catalogAttributes($entry, $existing->bundled));
        }

        $translation->update([
            'install_status' => TranslationInstallStatus::Pending,
            'install_step' => TranslationInstallStep::Queued,
            'install_error' => null,
        ]);

        Bus::dispatch(new InstallTranslationJob($translation->id));

        return $translation->fresh();
    }

    public function installBundled(string $moduleKey): Translation
    {
        $entry = $this->catalog->find($moduleKey);

        $translation = Translation::query()->firstOrCreate(
            ['abbrev' => strtolower($entry->short)],
            $this->attributesFromCatalog($entry, bundled: true),
        );

        if ($translation->isReady()) {
            return $translation;
        }

        $translation->update([
            'install_status' => TranslationInstallStatus::Pending,
            'install_step' => TranslationInstallStep::Queued,
            'install_error' => null,
            'bundled' => true,
        ]);

        Bus::dispatch(new InstallTranslationJob($translation->id));

        return $translation->fresh();
    }

    public function uninstall(string $moduleKey): void
    {
        $entry = $this->catalog->find($moduleKey);

        if ($this->registry->isBundled($entry->short)) {
            abort(422, "The {$entry->short} translation is bundled with the app and cannot be removed.");
        }

        $translation = Translation::query()->where('abbrev', strtolower($entry->short))->first();

        if ($translation === null || ! $translation->isReady()) {
            abort(404, "The {$entry->short} translation is not installed.");
        }

        app(TranslationSchemaManager::class)->dropTables($translation->abbrev);
        $translation->delete();
    }

    private function createTranslationRecord(CatalogEntry $entry): Translation
    {
        return Translation::query()->create($this->catalogAttributes($entry));
    }

    /** @return array<string, mixed> */
    private function catalogAttributes(CatalogEntry $entry, bool $bundled = false): array
    {
        return array_filter([
            'abbrev' => strtolower($entry->short),
            'name' => $entry->name,
            'format' => $entry->markupFormat?->value,
            'install_status' => TranslationInstallStatus::Pending,
            'install_step' => TranslationInstallStep::Pending,
            'bundled' => $bundled,
        ], fn($value) => $value !== null);
    }

    /** @return array<string, mixed> */
    private function attributesFromCatalog(CatalogEntry $entry, bool $bundled = false): array
    {
        return $this->catalogAttributes($entry, $bundled);
    }
}
