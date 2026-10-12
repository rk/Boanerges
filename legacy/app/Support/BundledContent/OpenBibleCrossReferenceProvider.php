<?php

namespace App\Support\BundledContent;

use App\Contracts\BundledContentProvider;
use App\Jobs\Bible\ImportCrossReferencesJob;
use App\Services\Bible\CrossReferenceService;
use Illuminate\Contracts\Queue\ShouldQueue;

class OpenBibleCrossReferenceProvider implements BundledContentProvider
{
    public function __construct(
        private CrossReferenceService $crossReferences,
    ) {}

    public function key(): string
    {
        return 'cross-references';
    }

    public function isImported(): bool
    {
        return $this->crossReferences->isImported();
    }

    public function requiresReadyTranslation(): bool
    {
        return false;
    }

    public function importJob(bool $force = false): ShouldQueue
    {
        return new ImportCrossReferencesJob(force: $force);
    }
}
