<?php

namespace App\Services\Study;

use Native\Desktop\DataObjects\Printer;
use Native\Desktop\Facades\System;

class StudyPrintService
{
    public function __construct(private StudyPrintHtmlBuilder $htmlBuilder) {}

    /**
     * @return list<array{name: string, displayName: string, description: string}>
     */
    public function printers(): array
    {
        if (! config('nativephp-internal.running', false)) {
            abort(503, 'Printing is only available in the desktop app.');
        }

        return array_map(
            fn(Printer $printer): array => [
                'name' => $printer->name,
                'displayName' => $printer->displayName,
                'description' => $printer->description,
            ],
            System::printers(),
        );
    }

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
    public function print(array $study, bool $includeUserWork, ?string $printerName = null): void
    {
        if (! config('nativephp-internal.running', false)) {
            abort(503, 'Printing is only available in the desktop app.');
        }

        $html = $this->htmlBuilder->build($study, $includeUserWork);
        $landscape = (int) $study['columnCount'] > 1;

        System::print($html, $this->resolvePrinter($printerName), [
            'landscape' => $landscape,
            'silent' => false,
            'usePrinterDefaultPageSize' => true,
        ]);
    }

    private function resolvePrinter(?string $printerName): ?Printer
    {
        if ($printerName === null || $printerName === '') {
            return null;
        }

        foreach (System::printers() as $printer) {
            if ($printer->name === $printerName) {
                return $printer;
            }
        }

        abort(422, 'Selected printer is not available.');
    }
}
