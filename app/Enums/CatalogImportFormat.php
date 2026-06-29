<?php

namespace App\Enums;

enum CatalogImportFormat: string
{
    case Sword = 'sword';
    case Usfm = 'usfm';
    case Accordance = 'accordance';
}
