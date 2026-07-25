<?php

use App\Http\Controllers\VerseListController;
use Illuminate\Support\Facades\Route;

Route::prefix('verse-lists')->name('verse-lists.')->group(function (): void {
    Route::get('/', [VerseListController::class, 'index'])->name('index');
    Route::get('/{id}', [VerseListController::class, 'show'])->name('show');
    Route::put('/', [VerseListController::class, 'store'])->name('store');
    Route::put('/{id}', [VerseListController::class, 'update'])->name('update');
    Route::delete('/{id}', [VerseListController::class, 'destroy'])->name('destroy');
});
