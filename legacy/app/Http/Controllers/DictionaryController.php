<?php

namespace App\Http\Controllers;

use App\Services\Dictionary\DictionaryService;
use App\Services\Dictionary\Webster1828Dictionary;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class DictionaryController extends Controller
{
    public function suggest(Request $request, Webster1828Dictionary $dictionary, DictionaryService $import): JsonResponse
    {
        if (! $import->isImported()) {
            return response()->json([
                'status' => 'importing',
                'message' => 'Dictionary is still being set up.',
            ], 503);
        }

        $query = (string) $request->query('q', '');
        $limit = min(max((int) $request->query('limit', 10), 1), 25);

        return response()->json([
            'suggestions' => $dictionary->suggest($query, $limit),
        ]);
    }

    public function show(string $word, Webster1828Dictionary $dictionary, DictionaryService $import): JsonResponse
    {
        if (! $import->isImported()) {
            return response()->json([
                'status' => 'importing',
                'message' => 'Dictionary is still being set up.',
            ], 503);
        }

        return response()->json($dictionary->lookup($word));
    }
}
