<?php

namespace App\Listeners;

use App\Enums\TranslationInstallStatus;
use App\Jobs\Bible\ImportCrossReferencesJob;
use App\Jobs\Bible\InstallTranslationJob;
use App\Models\Translation;
use App\Services\Bible\CrossReferenceService;
use App\Services\Bible\TranslationInstaller;
use Illuminate\Support\Facades\Bus;

class EnsureBundledData
{
    /** @var list<TranslationInstallStatus> */
    private const ACTIVE_INSTALL_STATUSES = [
        TranslationInstallStatus::Pending,
        TranslationInstallStatus::Downloading,
        TranslationInstallStatus::CreatingSchema,
        TranslationInstallStatus::Importing,
        TranslationInstallStatus::Verifying,
        TranslationInstallStatus::Indexing,
    ];

    public function __construct(
        private TranslationInstaller $installer,
        private CrossReferenceService $crossReferences,
    ) {}

    public function handle(): void
    {
        if (app()->environment('testing') && ! config('boanerges.seed_bundled_in_tests', false)) {
            return;
        }

        if (! Translation::query()->where('install_status', TranslationInstallStatus::Ready)->exists()) {
            foreach (config('boanerges.bundled_modules', []) as $module) {
                $existing = Translation::query()->where('abbrev', strtolower($module))->first();

                if ($existing?->isReady()) {
                    continue;
                }

                if ($existing !== null && in_array($existing->install_status, self::ACTIVE_INSTALL_STATUSES, true)) {
                    if (config('queue.default') === 'sync') {
                        Bus::dispatchSync(new InstallTranslationJob($existing->id));
                    }

                    continue;
                }

                $translation = $this->installer->installBundled($module);

                if (config('queue.default') === 'sync') {
                    Bus::dispatchSync(new InstallTranslationJob($translation->id));
                }
            }
        }

        if (! $this->crossReferences->isImported() && ! app()->environment('testing')) {
            if (config('queue.default') === 'sync') {
                Bus::dispatchSync(new ImportCrossReferencesJob());
            } else {
                Bus::dispatch(new ImportCrossReferencesJob());
            }
        }
    }
}
