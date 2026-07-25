<?php

namespace App\Providers;

use App\Services\Bible\BibleModuleManager;
use App\Services\Bible\Import\AccordanceFormatImporter;
use App\Services\Bible\Import\SwordFormatImporter;
use App\Services\Bible\Import\TranslationImporterRegistry;
use App\Services\Bible\Import\UsfmFormatImporter;
use App\Services\Bible\Markup\VerseMarkupConverterFactory;
use App\Services\Bible\Markup\VerseTextFormatter;
use App\Support\BundledContent\OpenBibleCrossReferenceProvider;
use App\Support\BundledContent\Webster1828DictionaryProvider;
use App\Support\BundledContentRegistry;
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

        $this->app->singleton(VerseTextFormatter::class, fn(): VerseTextFormatter => VerseMarkupConverterFactory::defaultFormatter());

        $this->app->singleton(TranslationImporterRegistry::class, function ($app): TranslationImporterRegistry {
            $registry = new TranslationImporterRegistry();

            foreach ([
                SwordFormatImporter::class,
                UsfmFormatImporter::class,
                AccordanceFormatImporter::class,
            ] as $importerClass) {
                $registry->register($app->make($importerClass));
            }

            return $registry;
        });

        $this->app->singleton(BundledContentRegistry::class, function ($app): BundledContentRegistry {
            $registry = new BundledContentRegistry();

            foreach ([
                Webster1828DictionaryProvider::class,
                OpenBibleCrossReferenceProvider::class,
            ] as $providerClass) {
                $registry->register($app->make($providerClass));
            }

            return $registry;
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
