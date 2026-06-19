<?php

namespace App\Services\Bible\Import;

use App\Services\Bible\Markup\VerseTextFormatter;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Support\Facades\DB;
use ZipArchive;

class UsfmImporter
{
    private const string ID_MARKER = '\\id';

    private const string HEADING_MARKER = '\\h';

    private const string CHAPTER_MARKER = '\\c';

    private const string VERSE_MARKER = '\\v';

    public function __construct(
        private TranslationSchemaManager $schema,
        private VerseTextFormatter $verseTextFormatter,
    ) {}

    public function importFromZip(string $abbrev, string $zipPath): void
    {
        $extractDir = sys_get_temp_dir() . '/boanerges-usfm-' . uniqid();
        mkdir($extractDir, 0755, true);

        try {
            $zip = new ZipArchive();

            if ($zip->open($zipPath) !== true) {
                throw new \RuntimeException('Failed to open USFM archive.');
            }

            $zip->extractTo($extractDir);
            $zip->close();

            $usfmFiles = $this->findUsfmFiles($extractDir);

            if ($usfmFiles === []) {
                throw new \RuntimeException('No USFM file found in archive.');
            }

            $abbrev = $this->schema->validateAbbrev($abbrev);
            $booksTable = $this->schema->booksTable($abbrev);
            $versesTable = $this->schema->versesTable($abbrev);

            DB::table($booksTable)->delete();
            DB::table($versesTable)->delete();

            foreach ($usfmFiles as $usfmFile) {
                $this->importFromFile($abbrev, $usfmFile, clearTables: false);
            }
        } finally {
            $this->deleteDirectory($extractDir);
        }
    }

    public function importFromFile(string $abbrev, string $path, bool $clearTables = true): void
    {
        $abbrev = $this->schema->validateAbbrev($abbrev);
        $booksTable = $this->schema->booksTable($abbrev);
        $versesTable = $this->schema->versesTable($abbrev);

        if ($clearTables) {
            DB::table($booksTable)->delete();
            DB::table($versesTable)->delete();
        }

        $content = file_get_contents($path);

        if ($content === false) {
            throw new \RuntimeException('Failed to read USFM file.');
        }

        $this->importContent($abbrev, $content, $booksTable, $versesTable);
    }

    private function importContent(string $abbrev, string $content, string $booksTable, string $versesTable): void
    {
        $bookDbId = null;
        $pendingOsisId = null;
        $chapter = 0;
        $verse = 0;
        $verseText = '';
        $maxChapter = 0;

        $flushVerse = function () use (&$bookDbId, &$chapter, &$verse, &$verseText, $versesTable): void {
            if ($bookDbId === null || $chapter <= 0 || $verse <= 0 || trim($verseText) === '') {
                return;
            }

            DB::table($versesTable)->insert([
                'book_id' => $bookDbId,
                'chapter' => $chapter,
                'verse' => $verse,
                'text' => trim($verseText),
                'plain_text' => $this->plainTextFromUsfm(trim($verseText)),
            ]);
        };

        $flushBook = function () use (&$bookDbId, &$pendingOsisId, &$chapter, &$verse, &$verseText, &$maxChapter, $booksTable, $flushVerse): void {
            $flushVerse();

            if ($bookDbId !== null) {
                DB::table($booksTable)->where('id', $bookDbId)->update(['chapters' => $maxChapter]);
            }

            $bookDbId = null;
            $pendingOsisId = null;
            $chapter = 0;
            $verse = 0;
            $verseText = '';
            $maxChapter = 0;
        };

        $ensureBook = function () use (&$bookDbId, &$pendingOsisId, $booksTable): void {
            if ($bookDbId !== null || $pendingOsisId === null) {
                return;
            }

            $bookDbId = $this->insertBook($booksTable, $pendingOsisId, strtoupper($pendingOsisId));
        };

        foreach (preg_split('/\r\n|\r|\n/', $content) as $line) {
            $line = rtrim($line);

            if ($line === '') {
                continue;
            }

            if (preg_match($this->markerPattern(self::ID_MARKER) . '\s+(\S+)/u', $line, $matches)) {
                $flushBook();

                $osisId = $this->osisIdFromUsfmCode($matches[1]);

                if (! $this->isCanonBook($osisId)) {
                    continue;
                }

                $pendingOsisId = $osisId;

                continue;
            }

            if ($pendingOsisId === null) {
                continue;
            }

            if ($bookDbId === null && preg_match($this->markerPattern(self::HEADING_MARKER) . '\s+(.+)$/u', $line, $matches)) {
                $bookDbId = $this->insertBook($booksTable, $pendingOsisId, trim($matches[1]));

                continue;
            }

            if (preg_match($this->markerPattern(self::CHAPTER_MARKER) . '\s+(\d+)/u', $line, $matches)) {
                $ensureBook();
                $flushVerse();
                $chapter = (int) $matches[1];
                $verse = 0;
                $verseText = '';
                $maxChapter = max($maxChapter, $chapter);

                continue;
            }

            if ($bookDbId === null) {
                continue;
            }

            if (preg_match($this->markerPattern(self::VERSE_MARKER) . '\s+(\d+)(?:-\d+)?\s*(.*)$/u', $line, $matches)) {
                $flushVerse();
                $verse = (int) $matches[1];
                $verseText = $matches[2];

                continue;
            }

            if ($verse <= 0 || $chapter <= 0) {
                continue;
            }

            if (str_starts_with($line, '\\')) {
                if (preg_match('/^\\\\[a-z0-9*]+\s*(.*)$/iu', $line, $matches) && trim($matches[1]) !== '') {
                    $verseText = trim($verseText . ' ' . $matches[1]);
                }

                continue;
            }

            $verseText = trim($verseText . ' ' . $line);
        }

        $flushBook();
    }

    private function insertBook(string $booksTable, string $osisId, string $name): int
    {
        return DB::table($booksTable)->insertGetId([
            'name' => $name,
            'osis_id' => $osisId,
            'testament' => $this->guessTestament($osisId),
            'chapters' => 0,
        ]);
    }

    private function markerPattern(string $marker): string
    {
        return '/^' . preg_quote($marker, '/');
    }

    private function osisIdFromUsfmCode(string $code): string
    {
        return strtolower(trim($code));
    }

    private function isCanonBook(string $osisId): bool
    {
        foreach (BibleCanon::books() as $book) {
            if ($book['osis'] === $osisId) {
                return true;
            }
        }

        return false;
    }

    private function plainTextFromUsfm(string $text): string
    {
        $text = preg_replace('/\\\\f\s.*?\\\\f\*/su', '', $text) ?? $text;
        $text = preg_replace('/\\\\x\s.*?\\\\x\*/su', '', $text) ?? $text;
        $text = preg_replace('/\\\\[+]?[a-z0-9]+\*?\s?/iu', '', $text) ?? $text;

        return trim(preg_replace('/\s+/u', ' ', $text) ?? '');
    }

    private function guessTestament(string $osisId): string
    {
        $nt = ['mat', 'mrk', 'luk', 'jhn', 'act', 'rom', '1co', '2co', 'gal', 'eph', 'php', 'col', '1th', '2th', '1ti', '2ti', 'tit', 'phm', 'heb', 'jas', '1pe', '2pe', '1jn', '2jn', '3jn', 'jud', 'rev'];

        return in_array(strtolower($osisId), $nt, true) ? 'nt' : 'ot';
    }

    /** @return list<string> */
    private function findUsfmFiles(string $dir): array
    {
        $files = [];
        $iterator = new \RecursiveIteratorIterator(
            new \RecursiveDirectoryIterator($dir, \FilesystemIterator::SKIP_DOTS),
        );

        foreach ($iterator as $file) {
            if (! $file->isFile()) {
                continue;
            }

            $filename = $file->getFilename();

            if (str_starts_with($filename, '._') || ! preg_match('/\.usfm$/i', $filename)) {
                continue;
            }

            $files[] = $file->getPathname();
        }

        sort($files, SORT_NATURAL);

        return $files;
    }

    private function deleteDirectory(string $directory): void
    {
        if (! is_dir($directory)) {
            return;
        }

        foreach (scandir($directory) ?: [] as $item) {
            if ($item === '.' || $item === '..') {
                continue;
            }

            $path = $directory . DIRECTORY_SEPARATOR . $item;

            if (is_dir($path)) {
                $this->deleteDirectory($path);
            } else {
                @unlink($path);
            }
        }

        @rmdir($directory);
    }
}
