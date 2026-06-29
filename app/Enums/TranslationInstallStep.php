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
}
