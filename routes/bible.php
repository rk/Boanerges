<?php

use App\Http\Controllers\BibleController;
use App\Http\Controllers\DictionaryController;
use App\Http\Middleware\EnsureBundledDataMiddleware;
use Illuminate\Support\Facades\Route;

Route::prefix('bible')->name('bible.')->middleware(EnsureBundledDataMiddleware::class)->group(function (): void {
    Route::get('/translations', [BibleController::class, 'translations'])
        ->name('translations.index');

    Route::get('/translations/catalog', [BibleController::class, 'catalog'])
        ->name('translations.catalog');

    Route::post('/translations/{module}/install', [BibleController::class, 'install'])
        ->name('translations.install');

    Route::get('/translations/{module}/install-status', [BibleController::class, 'installStatus'])
        ->name('translations.install-status');

    Route::delete('/translations/{module}', [BibleController::class, 'uninstall'])
        ->name('translations.uninstall');

    Route::get('/search', [BibleController::class, 'search'])
        ->name('search');

    Route::get('/cross-references', [BibleController::class, 'crossReferences'])
        ->name('cross-references');

    Route::get('/dictionary/suggest', [DictionaryController::class, 'suggest'])
        ->name('dictionary.suggest');

    Route::get('/dictionary/{word}', [DictionaryController::class, 'show'])
        ->where('word', '.*')
        ->name('dictionary.show');

    Route::get('/translations/{translation}/books', [BibleController::class, 'books'])
        ->name('books.index');

    Route::get('/translations/{translation}/books/{book}/chapters/{chapter}', [BibleController::class, 'chapter'])
        ->whereNumber('chapter')
        ->name('chapters.show');
});
