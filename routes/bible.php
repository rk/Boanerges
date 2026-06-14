<?php

use App\Http\Controllers\BibleController;
use Illuminate\Support\Facades\Route;

Route::prefix('bible')->name('bible.')->group(function (): void {
    Route::get('/translations', [BibleController::class, 'translations'])
        ->name('translations.index');

    Route::get('/translations/catalog', [BibleController::class, 'catalog'])
        ->name('translations.catalog');

    Route::post('/translations/{module}/install', [BibleController::class, 'install'])
        ->name('translations.install');

    Route::delete('/translations/{module}', [BibleController::class, 'uninstall'])
        ->name('translations.uninstall');

    Route::get('/translations/{translation}/books', [BibleController::class, 'books'])
        ->name('books.index');

    Route::get('/translations/{translation}/books/{book}/chapters/{chapter}', [BibleController::class, 'chapter'])
        ->whereNumber('chapter')
        ->name('chapters.show');
});
