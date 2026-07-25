<?php

namespace App\Listeners;

use App\Models\Translation;
use App\Services\Bible\TranslationInstaller;
use App\Support\BundledContentRegistry;
use Illuminate\Support\Facades\Bus;

class EnsureBundledData
{
    public function __construct(
        private TranslationInstaller $installer,
        private BundledContentRegistry $bundledContent,
    ) {}

    public function handle(): void
    {
        if (app()->environment('testing') && ! config('boanerges.seed_bundled_in_tests', false)) {
            return;
        }

        if (! Translation::query()->where('install_status', \App\Enums\TranslationInstallStatus::Ready)->exists()) {
            foreach (config('boanerges.bundled_modules', []) as $module) {
                $existing = Translation::query()->where('abbrev', strtolower($module))->first();

                if ($existing?->isReady()) {
                    continue;
                }

                if ($existing !== null && $existing->install_status->isActive()) {
                    if (config('queue.default') === 'sync') {
                        Bus::dispatchSync(new \App\Jobs\Bible\InstallTranslationJob($existing->id));
                    }

                    continue;
                }

                $translation = $this->installer->installBundled($module);

                if (config('queue.default') === 'sync') {
                    Bus::dispatchSync(new \App\Jobs\Bible\InstallTranslationJob($translation->id));
                }
            }
        }

        foreach ($this->bundledContent->pending() as $provider) {
            $job = $provider->importJob();

            if (config('queue.default') === 'sync') {
                Bus::dispatchSync($job);
            } else {
                Bus::dispatch($job);
            }
        }
    }
}
