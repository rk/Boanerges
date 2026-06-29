<?php

namespace App\Http\Controllers;

use App\Http\Requests\PrintStudyRequest;
use App\Services\Study\StudyPrintService;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Response;

class StudyPrintController extends Controller
{
    public function index(StudyPrintService $print): JsonResponse
    {
        return response()->json([
            'printers' => $print->printers(),
        ]);
    }

    public function store(PrintStudyRequest $request, StudyPrintService $print): JsonResponse|Response
    {
        $printerName = $request->validated('printerName');

        $path = $print->print(
            $request->studySettings(),
            (bool) $request->validated('includeUserWork'),
            is_string($printerName) ? $printerName : null,
        );

        if ($path !== null) {
            return response()->json([
                'path' => $path,
            ]);
        }

        return response()->noContent();
    }
}
