<?php

namespace App\Contracts;

use Illuminate\Contracts\Queue\ShouldQueue;

interface BundledContentProvider
{
    public function key(): string;

    public function isImported(): bool;

    public function requiresReadyTranslation(): bool;

    public function importJob(bool $force = false): ShouldQueue;
}
