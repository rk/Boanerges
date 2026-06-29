<?php

namespace App\Enums;

enum StudyColumnType: string
{
    case BibleSecondary = 'bible-secondary';
    case Notes = 'notes';
    case Scribe = 'scribe';
    case Search = 'search';
    case CrossReferences = 'cross-references';
}
