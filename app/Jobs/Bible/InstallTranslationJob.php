<?php

namespace App\Jobs\Bible;

use App\Models\Translation;
use App\Services\Bible\TranslationImportPipeline;
use App\Support\BundledContentRegistry;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;
use Illuminate\Support\Facades\Bus;

class InstallTranslationJob implements ShouldQueue
{
    use Queueable;

    public int $timeout = 600;

    public function __construct(
        public int $translationId,
    ) {}

    public function handle(
        TranslationImportPipeline $pipeline,
        BundledContentRegistry $bundledContent,
    ): void {
        $translation = Translation::query()->findOrFail($this->translationId);

        if ($translation->isReady()) {
            return;
        }

        $pipeline->run($translation);

        if ($translation->fresh()?->isReady() !== true) {
            return;
        }

        foreach ($bundledContent->pending() as $provider) {
            if (! $provider->requiresReadyTranslation()) {
                continue;
            }

            if (config('queue.default') === 'sync') {
                Bus::dispatchSync($provider->importJob());
            } else {
                Bus::dispatch($provider->importJob());
            }
        }
    }
}
