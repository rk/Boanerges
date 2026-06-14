<?php

namespace App\Http\Controllers;

use App\Exceptions\BibleModuleNotInstalledException;
use App\Http\Resources\Bible\BookResource;
use App\Http\Resources\Bible\CatalogTranslationResource;
use App\Http\Resources\Bible\ChapterResource;
use App\Http\Resources\Bible\TranslationResource;
use App\Services\Bible\BibleModuleManager;
use App\Services\Bible\BookCatalog;
use App\Services\Bible\ChapterReader;
use App\Services\Bible\InstalledTranslationRegistry;
use App\Services\Bible\TranslationCatalog;
use App\Services\Bible\TranslationInstaller;
use Illuminate\Http\JsonResponse;

class BibleController extends Controller
{
    public function translations(InstalledTranslationRegistry $registry): JsonResponse
    {
        return response()->json([
            'translations' => TranslationResource::collection($registry->all()),
        ]);
    }

    public function catalog(
        TranslationCatalog $catalog,
        BibleModuleManager $modules,
        InstalledTranslationRegistry $registry,
    ): JsonResponse {
        return response()->json([
            'translations' => $catalog->all()->map(
                fn($entry) => (new CatalogTranslationResource($entry, $modules, $registry))->resolve(),
            )->values(),
        ]);
    }

    public function install(
        string $module,
        TranslationInstaller $installer,
        InstalledTranslationRegistry $registry,
        BibleModuleManager $modules,
    ): JsonResponse {
        $installer->install($module);
        $modules->clearCache();

        return response()->json([
            'translation' => new TranslationResource($registry->findByModule($module)),
        ]);
    }

    public function uninstall(
        string $module,
        TranslationInstaller $installer,
    ): JsonResponse {
        $installer->uninstall($module);

        return response()->json([
            'message' => 'Translation removed.',
        ]);
    }

    public function books(
        string $translation,
        InstalledTranslationRegistry $registry,
        BibleModuleManager $modules,
        BookCatalog $catalog,
    ): JsonResponse {
        $config = $registry->find($translation);

        try {
            $bible = $modules->open($config->module);
        } catch (BibleModuleNotInstalledException $exception) {
            abort(503, $exception->getMessage());
        }

        return response()->json([
            'books' => BookResource::collection($catalog->books($bible)),
        ]);
    }

    public function chapter(
        string $translation,
        string $book,
        int $chapter,
        InstalledTranslationRegistry $registry,
        BibleModuleManager $modules,
        ChapterReader $reader,
    ): JsonResponse {
        $config = $registry->find($translation);

        try {
            $bible = $modules->open($config->module);
        } catch (BibleModuleNotInstalledException $exception) {
            abort(503, $exception->getMessage());
        }

        return response()->json([
            'chapter' => new ChapterResource($reader->read($bible, $book, $chapter)),
        ]);
    }
}
