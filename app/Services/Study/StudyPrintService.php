<?php

namespace App\Services\Study;

use App\Services\Bible\OsisBookId;
use Illuminate\Support\Facades\Storage;
use Native\Desktop\DataObjects\Printer;
use Native\Desktop\Dialog;
use Native\Desktop\Facades\System;

class StudyPrintService
{
    public const PDF_DESTINATION = '__pdf__';

    public function __construct(
        private StudyPrintHtmlBuilder $htmlBuilder,
    ) {}

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
    public function print(array $study, bool $includeUserWork, ?string $printerName = null): ?string
    {
        if (! config('nativephp-internal.running', false)) {
            abort(503, 'Printing is only available in the desktop app.');
        }

        $html = $this->htmlBuilder->build($study, $includeUserWork);
        $landscape = (int) $study['columnCount'] > 1;

        if ($printerName === self::PDF_DESTINATION) {
            return $this->exportPdf($html, $study, $landscape);
        }

        System::print($html, $this->resolvePrinter($printerName), [
            'landscape' => $landscape,
            'silent' => false,
            'usePrinterDefaultPageSize' => true,
            'printBackground' => true,
        ]);

        return null;
    }

    /**
     * @param  array{bookId: string, chapter: int}  $study
     */
    private function exportPdf(string $html, array $study, bool $landscape): ?string
    {
        $pdf = System::printToPDF($html, [
            'landscape' => $landscape,
            'printBackground' => true,
        ]);

        $filename = sprintf(
            '%s %d study.pdf',
            OsisBookId::displayName($study['bookId']),
            $study['chapter'],
        );

        $path = Dialog::new()
            ->title('Save study PDF')
            ->filter('PDF', ['pdf'])
            ->defaultPath($this->defaultPdfPath($filename))
            ->save();

        if ($path === null || $path === '') {
            return null;
        }

        file_put_contents($path, base64_decode($pdf, true));

        return $path;
    }

    private function defaultPdfPath(string $filename): string
    {
        if (config('nativephp-internal.running', false)) {
            return Storage::disk('downloads')->path($filename);
        }

        return $filename;
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
