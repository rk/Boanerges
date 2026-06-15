<?php

namespace App\Providers;

use App\Services\Bible\BibleModuleManager;
use App\Services\Bible\Markup\PairTagVerseMarkupConverter;
use App\Services\Bible\Markup\VerseTextFormatter;
use Carbon\CarbonImmutable;
use Illuminate\Support\Facades\Date;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\ServiceProvider;
use Illuminate\Validation\Rules\Password;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        $this->app->singleton(BibleModuleManager::class);

        $this->app->singleton(VerseTextFormatter::class, function (): VerseTextFormatter {
            return new VerseTextFormatter([
                // GBF italic (e.g. YLT supplied words)
                new PairTagVerseMarkupConverter('FI', 'Fi', 'em'),
            ]);
        });
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        $this->configureDefaults();
    }

    /**
     * Configure default behaviors for production-ready applications.
     */
    protected function configureDefaults(): void
    {
        Date::use(CarbonImmutable::class);

        DB::prohibitDestructiveCommands(
            app()->isProduction(),
        );

        Password::defaults(
            fn(): ?Password => app()->isProduction()
            ? Password::min(12)
                ->mixedCase()
                ->letters()
                ->numbers()
                ->symbols()
                ->uncompromised()
            : null,
        );
    }
}
