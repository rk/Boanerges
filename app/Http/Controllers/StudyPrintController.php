<?php

namespace App\Http\Controllers;

use App\Http\Requests\PrintStudyRequest;
use App\Services\Study\StudyPrintService;
use Illuminate\Http\Response;

class StudyPrintController extends Controller
{
    public function store(PrintStudyRequest $request, StudyPrintService $print): Response
    {
        $print->print(
            $request->studySettings(),
            (bool) $request->validated('includeUserWork'),
        );

        return response()->noContent();
    }
}
