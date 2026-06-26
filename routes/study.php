<?php

use App\Http\Controllers\StudyPrintController;
use Illuminate\Support\Facades\Route;

Route::prefix('study')->name('study.')->group(function (): void {
    Route::post('/print', [StudyPrintController::class, 'store'])
        ->name('print');
});
