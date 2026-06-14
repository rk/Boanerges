<?php

namespace App\Http\Controllers;

use App\Http\Requests\UpdateScribeChapterRequest;
use App\Services\Scribe\ScribeChapterStore;
use Illuminate\Http\JsonResponse;
use InvalidArgumentException;

class ScribeController extends Controller
{
    public function show(
        string $book,
        int $chapter,
        ScribeChapterStore $store,
    ): JsonResponse {
        try {
            return response()->json([
                'verses' => $store->get($book, $chapter),
            ]);
        } catch (InvalidArgumentException) {
            abort(422, 'Invalid book id.');
        }
    }

    public function update(
        string $book,
        int $chapter,
        UpdateScribeChapterRequest $request,
        ScribeChapterStore $store,
    ): JsonResponse {
        try {
            $store->put($book, $chapter, $request->validated('verses'));
        } catch (InvalidArgumentException) {
            abort(422, 'Invalid book id.');
        }

        return response()->json([
            'verses' => $store->get($book, $chapter),
        ]);
    }
}
