<?php

namespace App\Http\Controllers;

use App\Http\Requests\UpdateNotesChapterRequest;
use App\Services\Notes\NotesChapterStore;
use Illuminate\Http\JsonResponse;
use InvalidArgumentException;

class NotesController extends Controller
{
    public function show(
        string $book,
        int $chapter,
        NotesChapterStore $store,
    ): JsonResponse {
        try {
            return response()->json([
                'content' => $store->get($book, $chapter),
            ]);
        } catch (InvalidArgumentException) {
            abort(422, 'Invalid book id.');
        }
    }

    public function update(
        string $book,
        int $chapter,
        UpdateNotesChapterRequest $request,
        NotesChapterStore $store,
    ): JsonResponse {
        try {
            $store->put($book, $chapter, $request->validated('content'));
        } catch (InvalidArgumentException) {
            abort(422, 'Invalid book id.');
        }

        return response()->json([
            'content' => $store->get($book, $chapter),
        ]);
    }
}
