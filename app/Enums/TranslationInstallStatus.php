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
}
