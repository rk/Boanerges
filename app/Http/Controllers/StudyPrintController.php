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

    public function store(PrintStudyRequest $request, StudyPrintService $print): Response
    {
        $printerName = $request->validated('printerName');

        $print->print(
            $request->studySettings(),
            (bool) $request->validated('includeUserWork'),
            is_string($printerName) ? $printerName : null,
        );

        return response()->noContent();
    }
}
