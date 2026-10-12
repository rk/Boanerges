<?php

use App\Http\Controllers\NotesController;
use Illuminate\Support\Facades\Route;

Route::prefix('notes')->name('notes.')->group(function (): void {
    Route::get('/{book}/chapters/{chapter}', [NotesController::class, 'show'])
        ->whereNumber('chapter')
        ->name('chapters.show');

    Route::put('/{book}/chapters/{chapter}', [NotesController::class, 'update'])
        ->whereNumber('chapter')
        ->name('chapters.update');
});
