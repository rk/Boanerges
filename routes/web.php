<?php

use App\Http\Controllers\SettingsController;
use Illuminate\Support\Facades\Route;

Route::inertia('/', 'Study')->name('home');

Route::get('/settings/readability', [SettingsController::class, 'showReadability'])
    ->name('settings.readability.show');

Route::patch('/settings/readability', [SettingsController::class, 'updateReadability'])
    ->name('settings.readability.update');

Route::get('/settings/study', [SettingsController::class, 'showStudy'])
    ->name('settings.study.show');

Route::patch('/settings/study', [SettingsController::class, 'updateStudy'])
    ->name('settings.study.update');
