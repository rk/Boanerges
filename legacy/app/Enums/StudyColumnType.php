<?php

namespace App\Enums;

enum StudyColumnType: string
{
    case BibleSecondary = 'bible-secondary';
    case Notes = 'notes';
    case Scribe = 'scribe';
    case Search = 'search';
    case CrossReferences = 'cross-references';
    case Comparison = 'comparison';
    case VerseList = 'verse-list';
    case Dictionary = 'dictionary';

    public function label(): string
    {
        return match ($this) {
            self::BibleSecondary => 'Translation',
            self::Notes => 'Notes',
            self::Scribe => 'Scribe',
            self::Search => 'Search',
            self::CrossReferences => 'Cross References',
            self::Dictionary => 'Dictionary',
            self::Comparison => 'Comparison',
            self::VerseList => 'Verse List',
        };
    }

    public function menuLabel(): string
    {
        return match ($this) {
            self::CrossReferences => 'Cross-References',
            default => $this->label(),
        };
    }

    public function allowsDuplicate(): bool
    {
        return $this === self::BibleSecondary;
    }

    public function isPrintable(): bool
    {
        return in_array($this, [
            self::BibleSecondary,
            self::Notes,
            self::Scribe,
            self::VerseList,
        ], true);
    }

    public function menuId(): ?string
    {
        return match ($this) {
            self::Search => 'study.search',
            self::CrossReferences => 'study.cross-references',
            self::Dictionary => 'study.dictionary',
            self::Comparison => 'study.comparison',
            self::VerseList => 'study.verse-list',
            default => null,
        };
    }

    public function menuHotkey(): ?string
    {
        return match ($this) {
            self::Search => 'CmdOrCtrl+F',
            self::CrossReferences => 'CmdOrCtrl+Shift+R',
            self::Dictionary => 'CmdOrCtrl+Shift+D',
            default => null,
        };
    }
}
