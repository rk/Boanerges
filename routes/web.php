<?php

use App\Enums\TranslationInstallStatus;
use App\Http\Controllers\SettingsController;
use App\Models\Translation;
use Illuminate\Support\Facades\Route;
use Illuminate\Support\Facades\Schema;

Route::get('/', function () {
    if (! Schema::hasTable('translations')) {
        return inertia('Welcome');
    }

    $ready = Translation::query()->where('install_status', TranslationInstallStatus::Ready)->exists();

    return $ready
        ? inertia('Study')
        : inertia('Welcome');
})->name('home');

Route::get('/settings/readability', [SettingsController::class, 'showReadability'])
    ->name('settings.readability.show');

Route::patch('/settings/readability', [SettingsController::class, 'updateReadability'])
    ->name('settings.readability.update');

Route::get('/settings/study', [SettingsController::class, 'showStudy'])
    ->name('settings.study.show');

Route::patch('/settings/study', [SettingsController::class, 'updateStudy'])
    ->name('settings.study.update');
