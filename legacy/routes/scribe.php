<?php

use App\Http\Controllers\ScribeController;
use Illuminate\Support\Facades\Route;

Route::prefix('scribe')->name('scribe.')->group(function (): void {
    Route::get('/{book}/chapters/{chapter}', [ScribeController::class, 'show'])
        ->whereNumber('chapter')
        ->name('chapters.show');

    Route::put('/{book}/chapters/{chapter}', [ScribeController::class, 'update'])
        ->whereNumber('chapter')
        ->name('chapters.update');
});
