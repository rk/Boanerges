<?php

namespace App\Enums;

enum TranslationInstallStep: string
{
    case Pending = 'pending';
    case Queued = 'queued';
    case Downloading = 'downloading';
    case SourceReady = 'source_ready';
    case Downloaded = 'downloaded';
    case CreatingSchema = 'creating_schema';
    case Importing = 'importing';
    case Verifying = 'verifying';
    case Indexing = 'indexing';
    case Indexed = 'indexed';
    case Ready = 'ready';
    case Failed = 'failed';

    public function progressPercent(): int
    {
        return match ($this) {
            self::Ready => 100,
            self::Indexed => 95,
            self::Indexing => 85,
            self::Verifying => 75,
            self::Importing => 50,
            self::CreatingSchema => 30,
            self::Downloaded, self::SourceReady => 20,
            self::Downloading => 10,
            self::Failed => 0,
            default => 0,
        };
    }
}
