<?php

namespace App\Jobs\Bible;

use App\Models\Translation;
use App\Services\Bible\TranslationImportPipeline;
use App\Services\Dictionary\DictionaryBootstrap;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;

class InstallTranslationJob implements ShouldQueue
{
    use Queueable;

    public int $timeout = 600;

    public function __construct(
        public int $translationId,
    ) {}

    public function handle(
        TranslationImportPipeline $pipeline,
        DictionaryBootstrap $dictionaryBootstrap,
    ): void {
        $translation = Translation::query()->findOrFail($this->translationId);

        if ($translation->isReady()) {
            return;
        }

        $pipeline->run($translation);

        if ($translation->fresh()?->isReady() === true) {
            $dictionaryBootstrap->dispatchIfNeeded();
        }
    }
}
