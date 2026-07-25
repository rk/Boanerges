<?php

namespace App\Support\BundledContent;

use App\Contracts\BundledContentProvider;
use App\Jobs\Bible\ImportWebster1828DictionaryJob;
use App\Services\Dictionary\DictionaryService;
use Illuminate\Contracts\Queue\ShouldQueue;

class Webster1828DictionaryProvider implements BundledContentProvider
{
    public function __construct(
        private DictionaryService $dictionary,
    ) {}

    public function key(): string
    {
        return 'dictionary';
    }

    public function isImported(): bool
    {
        return $this->dictionary->isImported();
    }

    public function requiresReadyTranslation(): bool
    {
        return true;
    }

    public function importJob(bool $force = false): ShouldQueue
    {
        return new ImportWebster1828DictionaryJob(force: $force);
    }
}
