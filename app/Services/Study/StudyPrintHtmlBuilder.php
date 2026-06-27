<?php

namespace App\Services\Study;

use App\Services\Bible\DbChapterReader;
use App\Services\Bible\InstalledTranslationRegistry;
use App\Services\Notes\NotesChapterStore;
use App\Services\ReadabilitySettingsStore;
use InvalidArgumentException;

class StudyPrintHtmlBuilder
{
    private const COLUMN_LABELS = [
        'bible-secondary' => 'Translation',
        'notes' => 'Notes',
        'scribe' => 'Scribe',
        'search' => 'Search',
        'cross-references' => 'Cross References',
    ];

    public function __construct(
        private DbChapterReader $chapters,
        private NotesChapterStore $notes,
        private InstalledTranslationRegistry $translations,
        private ReadabilitySettingsStore $readability,
    ) {}

    public function build(array $study, bool $includeUserWork): string
    {
        $bookId = $study['bookId'];
        $chapterNumber = (int) $study['chapter'];
        $primaryChapter = $this->chapters->read($study['translationId'], $bookId, $chapterNumber);
        $readability = $this->readability->get();
        $columns = $this->buildColumns($study, $bookId, $chapterNumber, $primaryChapter, $includeUserWork);

        return view('print.study', [
            'landscape' => (int) $study['columnCount'] > 1,
            'fontFamily' => $this->printFontFamily((string) $readability['fontFamily']),
            'fontSize' => (int) $readability['fontSize'],
            'lineHeight' => (float) $readability['lineHeight'],
            'justifyText' => (bool) $readability['justifyText'],
            'columns' => $columns,
        ])->render();
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
     * @param  array{
     *     book: string,
     *     bookAbbrev: string,
     *     chapter: int,
     *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
     * }  $primaryChapter
     * @return list<array{
     *     label: string,
     *     kind: string,
     *     verses?: list<array{number: int, text: string, paragraphStart?: bool}>,
     *     content?: string,
     *     linedVerses?: list<array{number: int, paragraphStart?: bool}>,
     *     message?: string
     * }>
     */
    private function buildColumns(
        array $study,
        string $bookId,
        int $chapterNumber,
        array $primaryChapter,
        bool $includeUserWork,
    ): array {
        $primaryTranslation = $this->translations->find($study['translationId']);
        $columns = [[
            'label' => sprintf(
                '%s %d (%s)',
                $primaryChapter['book'],
                $chapterNumber,
                strtoupper($primaryTranslation->abbrev),
            ),
            'kind' => 'bible',
            'verses' => $primaryChapter['verses'],
        ]];

        foreach ($study['columns'] as $slotIndex => $type) {
            $columns[] = match ($type) {
                'bible-secondary' => $this->bibleColumn(
                    $study,
                    $bookId,
                    $chapterNumber,
                    $slotIndex === 0 ? $study['translationBId'] : $study['translationCId'],
                ),
                'notes' => $this->notesColumn($bookId, $chapterNumber, $primaryChapter['book'], $includeUserWork),
                'scribe' => $this->scribeColumn($primaryChapter),
                'search', 'cross-references' => [
                    'label' => self::COLUMN_LABELS[$type],
                    'kind' => 'message',
                    'message' => 'Interactive view — not included in print.',
                ],
                default => throw new InvalidArgumentException("Unknown column type: {$type}"),
            };
        }

        return $columns;
    }

    /**
     * @return array{
     *     label: string,
     *     kind: string,
     *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
     * }
     */
    private function bibleColumn(array $study, string $bookId, int $chapterNumber, string $translationId): array
    {
        $translation = $this->translations->find($translationId);
        $chapter = $this->chapters->read($translationId, $bookId, $chapterNumber);

        return [
            'label' => sprintf(
                '%s %d (%s)',
                $chapter['book'],
                $chapterNumber,
                strtoupper($translation->abbrev),
            ),
            'kind' => 'bible',
            'verses' => $chapter['verses'],
        ];
    }

    /**
     * @return array{label: string, kind: string, content?: string, linedVerses?: list<array{number: int, paragraphStart?: bool}>}
     */
    private function notesColumn(
        string $bookId,
        int $chapterNumber,
        string $bookName,
        bool $includeUserWork,
    ): array {
        $label = sprintf('%s %d', $bookName, $chapterNumber);

        if ($includeUserWork) {
            return [
                'label' => $label,
                'kind' => 'notes',
                'content' => $this->notes->get($bookId, $chapterNumber),
            ];
        }

        return [
            'label' => $label,
            'kind' => 'lined-notes',
        ];
    }

    /**
     * @param  array{
     *     book: string,
     *     bookAbbrev: string,
     *     chapter: int,
     *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
     * }  $primaryChapter
     * @return array{
     *     label: string,
     *     kind: string,
     *     linedVerses: list<array{number: int, paragraphStart?: bool}>
     * }
     */
    private function scribeColumn(array $primaryChapter): array
    {
        $linedVerses = array_map(
            fn(array $verse): array => [
                'number' => $verse['number'],
                'paragraphStart' => $verse['paragraphStart'] ?? false,
            ],
            $primaryChapter['verses'],
        );

        return [
            'label' => sprintf('%s %d', $primaryChapter['book'], $primaryChapter['chapter']),
            'kind' => 'scribe',
            'linedVerses' => $linedVerses,
        ];
    }

    private function printFontFamily(string $fontFamily): string
    {
        return match ($fontFamily) {
            'sans-serif' => 'system-ui, sans-serif',
            default => 'Georgia, "Times New Roman", serif',
        };
    }
}
