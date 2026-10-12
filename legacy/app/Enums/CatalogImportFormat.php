<?php

namespace App\Enums;

enum CatalogImportFormat: string
{
    case Sword = 'sword';
    case Usfm = 'usfm';
    case Accordance = 'accordance';

    public function defaultMarkupFormat(): ?VerseMarkupFormat
    {
        return match ($this) {
            self::Usfm => VerseMarkupFormat::Usfm,
            default => null,
        };
    }

    public function usesSwordConf(): bool
    {
        return $this === self::Sword;
    }
}
