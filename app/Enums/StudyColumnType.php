<?php

namespace App\Enums;

enum StudyColumnType: string
{
    case BibleSecondary = 'bible-secondary';
    case Notes = 'notes';
    case Scribe = 'scribe';
    case Search = 'search';
    case CrossReferences = 'cross-references';
    case Dictionary = 'dictionary';
    case Comparison = 'comparison';
    case VerseList = 'verse-list';
}
