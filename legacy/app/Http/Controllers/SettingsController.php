<?php

namespace App\Http\Controllers;

use App\Http\Requests\UpdateReadabilitySettingsRequest;
use App\Http\Requests\UpdateStudySettingsRequest;
use App\Services\ReadabilitySettingsStore;
use App\Services\StudySettingsStore;
use App\Support\ApplicationMenuBuilder;
use Illuminate\Http\JsonResponse;

class SettingsController extends Controller
{
    public function showReadability(ReadabilitySettingsStore $settings): JsonResponse
    {
        return response()->json([
            'readability' => $settings->get(),
        ]);
    }

    public function updateReadability(
        UpdateReadabilitySettingsRequest $request,
        ReadabilitySettingsStore $settings,
    ): JsonResponse {
        return response()->json([
            'readability' => $settings->update($request->validated()),
        ]);
    }

    public function showStudy(StudySettingsStore $settings): JsonResponse
    {
        return response()->json([
            'study' => $settings->get(),
        ]);
    }

    public function updateStudy(
        UpdateStudySettingsRequest $request,
        StudySettingsStore $settings,
        ApplicationMenuBuilder $menu,
    ): JsonResponse {
        $study = $settings->update($request->validated());
        $menu->register();

        return response()->json([
            'study' => $study,
        ]);
    }
}
