<?php

namespace App\Enums;

enum TranslationInstallStatus: string
{
    case Pending = 'pending';
    case Downloading = 'downloading';
    case CreatingSchema = 'creating_schema';
    case Importing = 'importing';
    case Verifying = 'verifying';
    case Indexing = 'indexing';
    case Ready = 'ready';
    case Failed = 'failed';

    public function progressPercent(): ?int
    {
        return match ($this) {
            self::Ready => 100,
            self::Indexing => 85,
            self::Verifying => 75,
            self::Importing => 50,
            self::CreatingSchema => 30,
            self::Downloading => 10,
            self::Failed => 0,
            default => null,
        };
    }

    public function isActive(): bool
    {
        return match ($this) {
            self::Pending,
            self::Downloading,
            self::CreatingSchema,
            self::Importing,
            self::Verifying,
            self::Indexing => true,
            default => false,
        };
    }
}
