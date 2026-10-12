<?php

namespace App\Http\Controllers;

use App\Http\Requests\StoreVerseListRequest;
use App\Services\VerseList\VerseListStore;
use Illuminate\Http\JsonResponse;
use InvalidArgumentException;

class VerseListController extends Controller
{
    public function index(VerseListStore $store): JsonResponse
    {
        return response()->json([
            'lists' => $store->list(),
        ]);
    }

    public function show(string $id, VerseListStore $store): JsonResponse
    {
        try {
            return response()->json([
                'list' => $store->get($id),
            ]);
        } catch (InvalidArgumentException) {
            abort(404, 'Verse list not found.');
        }
    }

    public function store(
        StoreVerseListRequest $request,
        VerseListStore $store,
    ): JsonResponse {
        return $this->saveList($request, $store);
    }

    public function update(
        StoreVerseListRequest $request,
        VerseListStore $store,
        string $id,
    ): JsonResponse {
        return $this->saveList($request, $store, $id);
    }

    public function destroy(string $id, VerseListStore $store): JsonResponse
    {
        try {
            $store->delete($id);
        } catch (InvalidArgumentException) {
            abort(404, 'Verse list not found.');
        }

        return response()->json([
            'deleted' => true,
        ]);
    }

    private function saveList(
        StoreVerseListRequest $request,
        VerseListStore $store,
        ?string $id = null,
    ): JsonResponse {
        try {
            $list = $store->save(
                $request->validated('title'),
                $request->validated('entries'),
                $id,
            );
        } catch (InvalidArgumentException $exception) {
            abort(422, $exception->getMessage());
        }

        return response()->json([
            'list' => $list,
        ]);
    }
}
