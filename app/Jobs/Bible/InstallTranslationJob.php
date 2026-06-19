<?php

namespace App\Jobs\Bible;

use App\Models\Translation;
use App\Services\Bible\TranslationImportPipeline;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Queue\Queueable;

class InstallTranslationJob implements ShouldQueue
{
    use Queueable;

    public int $timeout = 600;

    public function __construct(
        public int $translationId,
    ) {}

    public function handle(TranslationImportPipeline $pipeline): void
    {
        $translation = Translation::query()->findOrFail($this->translationId);

        if ($translation->isReady()) {
            return;
        }

        $pipeline->run($translation);
    }
}
