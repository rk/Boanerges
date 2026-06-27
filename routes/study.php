<?php

use App\Http\Controllers\StudyPrintController;
use Illuminate\Support\Facades\Route;

Route::prefix('study')->name('study.')->group(function (): void {
    Route::get('/printers', [StudyPrintController::class, 'index'])
        ->name('printers.index');

    Route::post('/print', [StudyPrintController::class, 'store'])
        ->name('print');
});
