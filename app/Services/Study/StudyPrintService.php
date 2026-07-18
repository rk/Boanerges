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

    public const HTML_DESTINATION = '__html__';

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

        $printers = array_map(
            fn(Printer $printer): array => [
                'name' => $printer->name,
                'displayName' => $printer->displayName,
                'description' => $printer->description,
            ],
            System::printers(),
        );

        /** @var list<array{name: string, displayName: string, description: string}> $printers */
        return $printers;
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

        if ($printerName === self::HTML_DESTINATION) {
            return $this->exportHtml($html, $study);
        }

        System::print($html, $this->resolvePrinter($printerName), $this->printSettings($landscape));

        return null;
    }

    /**
     * @return array{
     *     landscape: bool,
     *     pageSize: string,
     *     silent: bool,
     *     usePrinterDefaultPageSize: bool,
     *     printBackground: bool
     * }
     */
    private function printSettings(bool $landscape): array
    {
        return [
            'landscape' => $landscape,
            'pageSize' => 'A4',
            'silent' => false,
            'usePrinterDefaultPageSize' => false,
            'printBackground' => true,
        ];
    }

    /**
     * @return array{landscape: bool, pageSize: string, printBackground: bool}
     */
    private function pdfSettings(bool $landscape): array
    {
        return [
            'landscape' => $landscape,
            'pageSize' => 'A4',
            'printBackground' => true,
        ];
    }

    /**
     * @param  array{bookId: string, chapter: int}  $study
     */
    private function exportPdf(string $html, array $study, bool $landscape): ?string
    {
        $pdf = System::printToPDF($html, $this->pdfSettings($landscape));

        $filename = sprintf(
            '%s %d study.pdf',
            OsisBookId::displayName($study['bookId']),
            $study['chapter'],
        );

        $path = Dialog::new()
            ->title('Save study PDF')
            ->filter('PDF', ['pdf'])
            ->defaultPath($this->defaultExportPath($filename))
            ->save();

        if ($path === null || $path === '') {
            return null;
        }

        file_put_contents($path, base64_decode($pdf, true));

        return $path;
    }

    /**
     * @param  array{bookId: string, chapter: int}  $study
     */
    private function exportHtml(string $html, array $study): ?string
    {
        $filename = sprintf(
            '%s %d study.html',
            OsisBookId::displayName($study['bookId']),
            $study['chapter'],
        );

        $path = Dialog::new()
            ->title('Save study HTML')
            ->filter('HTML', ['html', 'htm'])
            ->defaultPath($this->defaultExportPath($filename))
            ->save();

        if ($path === null || $path === '') {
            return null;
        }

        file_put_contents($path, $html);

        return $path;
    }

    private function defaultExportPath(string $filename): string
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
