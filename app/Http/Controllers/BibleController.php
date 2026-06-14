<?php

namespace App\Http\Controllers;

use App\Exceptions\BibleModuleNotInstalledException;
use App\Http\Resources\Bible\BookResource;
use App\Http\Resources\Bible\ChapterResource;
use App\Http\Resources\Bible\TranslationResource;
use App\Services\Bible\BibleModuleManager;
use App\Services\Bible\BookCatalog;
use App\Services\Bible\ChapterReader;
use App\Services\Bible\InstalledTranslationRegistry;
use Illuminate\Http\JsonResponse;

class BibleController extends Controller
{
    public function translations(InstalledTranslationRegistry $registry): JsonResponse
    {
        return response()->json([
            'translations' => TranslationResource::collection($registry->all()),
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
