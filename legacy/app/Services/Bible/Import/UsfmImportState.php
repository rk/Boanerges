<?php

namespace App\Services\Bible\Import;

final class UsfmImportState
{
    public ?int $bookDbId = null;

    public ?string $pendingOsisId = null;

    public int $chapter = 0;

    public int $verse = 0;

    public string $verseText = '';

    public int $maxChapter = 0;
}
