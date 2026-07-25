<?php

namespace App\Services\Study;

use App\Services\Bible\DbChapterReader;
use App\Services\Bible\InstalledTranslationRegistry;
use App\Services\Bible\OsisBookId;
use App\Services\Notes\NotesChapterStore;
use App\Services\ReadabilitySettingsStore;
use App\Services\Scribe\ScribeChapterStore;
use InvalidArgumentException;

class StudyPrintHtmlBuilder
{
    private const COLUMN_LABELS = [
        'bible-secondary' => 'Translation',
        'notes' => 'Notes',
        'scribe' => 'Scribe',
        'search' => 'Search',
        'cross-references' => 'Cross References',
        'dictionary' => 'Dictionary',
        'comparison' => 'Comparison',
        'verse-list' => 'Verse List',
    ];

    public function __construct(
        private DbChapterReader $chapters,
        private NotesChapterStore $notes,
        private ScribeChapterStore $scribe,
        private InstalledTranslationRegistry $translations,
        private ReadabilitySettingsStore $readability,
    ) {}

    /**
     * @param  array{
     *     columnCount: int,
     *     columns: list<string>,
     *     bookId: string,
     *     chapter: int,
     *     translationId: string,
     *     translationBId: string,
     *     translationCId: string,
     *     verseList?: array{title: string, entries: list<array{bookId: string, chapter: int, verse: int}>}|null
     * }  $study
     */
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
     *     translationCId: string,
     *     verseList?: array{title: string, entries: list<array{bookId: string, chapter: int, verse: int}>}|null
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
                'scribe' => $this->scribeColumn($primaryChapter, $bookId, $chapterNumber, $includeUserWork),
                'search', 'cross-references', 'dictionary', 'comparison' => [
                    'label' => self::COLUMN_LABELS[$type],
                    'kind' => 'message',
                    'message' => 'Interactive view — not included in print.',
                ],
                'verse-list' => $this->verseListColumn($study),
                default => throw new InvalidArgumentException("Unknown column type: {$type}"),
            };
        }

        return $columns;
    }

    /**
     * @param  array{
     *     columnCount: int,
     *     columns: list<string>,
     *     bookId: string,
     *     chapter: int,
     *     translationId: string,
     *     translationBId: string,
     *     translationCId: string,
     *     verseList?: array{title: string, entries: list<array{bookId: string, chapter: int, verse: int}>}|null
     * }  $study
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
     * @return array{label: string, kind: string, content?: string, verses?: list<array{number: int, text: string, paragraphStart?: bool}>}
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
                'label' => $label . '(Notes)',
                'kind' => 'notes',
                'content' => $this->notes->get($bookId, $chapterNumber),
            ];
        }

        return [
            'label' => $label . '(Notes)',
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
     * @return array{label: string, kind: string, verses?: list<array{number: int, text: string, paragraphStart?: bool}>}
     */
    private function scribeColumn(
        array $primaryChapter,
        string $bookId,
        int $chapterNumber,
        bool $includeUserWork,
    ): array {
        $label = sprintf('%s %d', $primaryChapter['book'], $primaryChapter['chapter']);

        if (! $includeUserWork) {
            return [
                'label' => $label . '(Scribe)',
                'kind' => 'scribe',
            ];
        }

        return [
            'label' => $label . '(Scribe)',
            'kind' => 'scribe-content',
            'verses' => $this->scribeContentVerses($primaryChapter, $bookId, $chapterNumber),
        ];
    }

    /**
     * @param  array{
     *     verses: list<array{number: int, text: string, paragraphStart?: bool}>
     * }  $primaryChapter
     * @return list<array{number: int, text: string, paragraphStart?: bool}>
     */
    private function scribeContentVerses(array $primaryChapter, string $bookId, int $chapterNumber): array
    {
        $draft = collect($this->scribe->get($bookId, $chapterNumber))->keyBy('verse');

        $verses = [];
        $pendingBreak = false;

        foreach ($primaryChapter['verses'] as $source) {
            /** @var array{verse: int, text: string, paragraphStart?: bool}|null $entry */
            $entry = $draft->get($source['number']);
            $text = (string) ($entry['text'] ?? '');
            $hasText = trim($text) !== '';
            $startsParagraph = $pendingBreak || $this->effectiveScribeParagraphStart(
                $source['number'],
                $entry['paragraphStart'] ?? null,
            );
            $pendingBreak = false;

            if (! $hasText) {
                if ($startsParagraph && $verses !== []) {
                    $pendingBreak = true;
                }

                continue;
            }

            $verses[] = [
                'number' => $source['number'],
                'text' => $text,
                'paragraphStart' => $verses === [] ? true : $startsParagraph,
            ];
        }

        return $verses;
    }

    private function effectiveScribeParagraphStart(
        int $verseNumber,
        ?bool $override,
    ): bool {
        if ($override !== null) {
            return $override;
        }

        return $verseNumber === 1;
    }

    private function printFontFamily(string $fontFamily): string
    {
        return match ($fontFamily) {
            'sans-serif' => 'system-ui, sans-serif',
            default => 'Georgia, "Times New Roman", serif',
        };
    }

    /**
     * @param  array{
     *     translationId: string,
     *     verseList?: array{title: string, entries: list<array{bookId: string, chapter: int, verse: int}>}|null
     * }  $study
     * @return array{
     *     label: string,
     *     kind: string,
     *     entries: list<array{label: string, text: string}>
     * }
     */
    private function verseListColumn(array $study): array
    {
        $verseList = $study['verseList'] ?? null;

        if (! is_array($verseList)) {
            return [
                'label' => 'Verse List',
                'kind' => 'verse-list',
                'entries' => [],
            ];
        }

        $title = (string) $verseList['title'];
        $entries = $verseList['entries'];
        $translationId = (string) $study['translationId'];
        $printedEntries = [];

        foreach ($entries as $entry) {
            $bookId = (string) $entry['bookId'];
            $chapter = (int) $entry['chapter'];
            $verse = (int) $entry['verse'];

            if ($bookId === '' || $chapter < 1 || $verse < 1) {
                continue;
            }

            $text = '—';

            try {
                $chapterData = $this->chapters->read($translationId, $bookId, $chapter);

                foreach ($chapterData['verses'] as $verseData) {
                    if ((int) $verseData['number'] === $verse) {
                        $text = (string) $verseData['text'];
                        break;
                    }
                }
            } catch (\Throwable) {
                $text = '—';
            }

            $printedEntries[] = [
                'label' => sprintf(
                    '%s %d:%d',
                    OsisBookId::displayName($bookId),
                    $chapter,
                    $verse,
                ),
                'text' => $text,
            ];
        }

        return [
            'label' => $title,
            'kind' => 'verse-list',
            'entries' => $printedEntries,
        ];
    }
}
