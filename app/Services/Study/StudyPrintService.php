<?php

namespace App\Services\Study;

use Native\Desktop\Facades\System;

class StudyPrintService
{
    public function __construct(private StudyPrintHtmlBuilder $htmlBuilder) {}

    /**
     * @param  array{
     *     columnCount: int,
     *     columns: list<string>,
     *     bookId: string,
     *     chapter: int,
     *     translationId: string,
     *     translationBId: string,
     *     translationCId: string
     * }  $study
     */
    public function print(array $study, bool $includeUserWork): void
    {
        if (! config('nativephp-internal.running', false)) {
            abort(503, 'Printing is only available in the desktop app.');
        }

        $html = $this->htmlBuilder->build($study, $includeUserWork);
        $landscape = (int) $study['columnCount'] > 1;

        System::print($html, null, [
            'pageSize' => 'A4',
            'landscape' => $landscape,
        ]);
    }
}
