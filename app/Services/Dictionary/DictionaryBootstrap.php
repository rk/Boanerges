<?php

namespace App\Services\Dictionary;

use App\Enums\TranslationInstallStatus;
use App\Jobs\Bible\ImportWebster1828DictionaryJob;
use App\Models\Translation;
use Illuminate\Support\Facades\Bus;

class DictionaryBootstrap
{
    public function __construct(
        private DictionaryService $dictionary,
    ) {}

    public function dispatchIfNeeded(): void
    {
        if (app()->environment('testing') && ! config('boanerges.seed_bundled_in_tests', false)) {
            return;
        }

        if ($this->dictionary->isImported()) {
            return;
        }

        if (! Translation::query()->where('install_status', TranslationInstallStatus::Ready)->exists()) {
            return;
        }

        if (config('queue.default') === 'sync') {
            Bus::dispatchSync(new ImportWebster1828DictionaryJob());
        } else {
            Bus::dispatch(new ImportWebster1828DictionaryJob());
        }
    }
}
