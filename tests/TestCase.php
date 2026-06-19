<?php

namespace Tests;

use App\Listeners\EnsureBundledData;
use Illuminate\Foundation\Testing\TestCase as BaseTestCase;
use Illuminate\Support\Facades\Storage;

abstract class TestCase extends BaseTestCase
{
    protected function afterRefreshingDatabase()
    {
        if (is_dir(Storage::disk('extras')->path('sword/mods.d'))) {
            $this->app->make(EnsureBundledData::class)->handle();
        }
    }
}
