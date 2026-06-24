<?php

namespace App\Http\Controllers;

use App\Http\Resources\Bible\BookResource;
use App\Http\Resources\Bible\CatalogTranslationResource;
use App\Http\Resources\Bible\ChapterResource;
use App\Http\Resources\Bible\TranslationResource;
use App\Services\Bible\CrossReferenceService;
use App\Services\Bible\DbBookCatalog;
use App\Services\Bible\DbChapterReader;
use App\Services\Bible\InstalledTranslationRegistry;
use App\Services\Bible\SearchService;
use App\Services\Bible\TranslationCatalog;
use App\Services\Bible\TranslationInstaller;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

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
        InstalledTranslationRegistry $registry,
    ): JsonResponse {
        return response()->json([
            'translations' => $catalog->all()->map(
                fn($entry) => (new CatalogTranslationResource($entry, $registry))->resolve(),
            )->values(),
        ]);
    }

    public function install(
        string $module,
        TranslationInstaller $installer,
        InstalledTranslationRegistry $registry,
    ): JsonResponse {
        $translation = $installer->install($module);

        return response()->json([
            'translation' => new TranslationResource($registry->toConfig($translation)),
        ], 202);
    }

    public function installStatus(
        string $module,
        InstalledTranslationRegistry $registry,
    ): JsonResponse {
        $model = $registry->findModel($module);

        if ($model === null) {
            abort(404, 'Translation not found.');
        }

        return response()->json([
            'abbrev' => $model->abbrev,
            'install_status' => $model->install_status->value,
            'step' => $model->install_step ?? 'pending',
            'percent' => $this->percentForStatus($model->install_status->value, $model->install_step),
            'install_error' => $model->install_error,
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
        DbBookCatalog $catalog,
    ): JsonResponse {
        $config = $registry->find($translation);

        return response()->json([
            'books' => BookResource::collection($catalog->books($config->id)),
        ]);
    }

    public function chapter(
        string $translation,
        string $book,
        int $chapter,
        InstalledTranslationRegistry $registry,
        DbChapterReader $reader,
    ): JsonResponse {
        $config = $registry->find($translation);

        return response()->json([
            'chapter' => new ChapterResource($reader->read($config->id, $book, $chapter)),
        ]);
    }

    public function search(Request $request, SearchService $search): JsonResponse
    {
        $validated = $request->validate([
            'q' => ['required', 'string', 'min:2'],
            'translation' => ['nullable', 'string'],
            'limit' => ['nullable', 'integer', 'min:1', 'max:100'],
        ]);

        return response()->json([
            'results' => $search->search(
                query: $validated['q'],
                translation: $validated['translation'] ?? null,
                limit: $validated['limit'] ?? 50,
            ),
        ]);
    }

    public function crossReferences(
        Request $request,
        CrossReferenceService $crossReferences,
    ): JsonResponse {
        $validated = $request->validate([
            'book' => ['required', 'string'],
            'chapter' => ['required', 'integer', 'min:1'],
            'verse' => ['required', 'integer', 'min:1'],
        ]);

        return response()->json([
            'references' => $crossReferences->forVerse(
                $validated['book'],
                (int) $validated['chapter'],
                (int) $validated['verse'],
            ),
        ]);
    }

    private function percentForStatus(string $status, ?string $step): int
    {
        return match ($status) {
            'ready' => 100,
            'indexing' => 85,
            'verifying' => 75,
            'importing' => 50,
            'creating_schema' => 30,
            'downloading' => 10,
            'failed' => 0,
            default => match ($step) {
                'ready' => 100,
                'indexed' => 95,
                'indexing' => 85,
                'verifying' => 75,
                'importing' => 50,
                'creating_schema' => 30,
                'downloaded', 'source_ready' => 20,
                'downloading' => 10,
                default => 0,
            },
        };
    }
}
