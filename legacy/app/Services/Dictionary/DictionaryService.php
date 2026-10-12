<?php

namespace App\Services\Dictionary;

use App\Jobs\Bible\ImportWebster1828DictionaryJob;
use Illuminate\Support\Facades\DB;

class DictionaryService
{
    public function isImported(): bool
    {
        return DB::table('import_meta')
            ->where('key', ImportWebster1828DictionaryJob::META_KEY)
            ->whereNotNull('completed_at')
            ->exists();
    }
}
