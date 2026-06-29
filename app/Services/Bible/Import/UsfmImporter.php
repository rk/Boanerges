<?php

namespace App\Services\Bible\Import;

use App\Services\Bible\Markup\VerseTextFormatter;
use App\Services\Bible\OsisBookId;
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

        $this->importContent($content, $booksTable, $versesTable);
    }

    private function importContent(string $content, string $booksTable, string $versesTable): void
    {
        $state = new UsfmImportState();
        $lines = preg_split('/\r\n|\r|\n/', $content);

        if ($lines === false) {
            return;
        }

        foreach ($lines as $line) {
            $line = rtrim($line);

            if ($line === '') {
                continue;
            }

            if (preg_match($this->markerPattern(self::ID_MARKER) . '\s+(\S+)/u', $line, $matches)) {
                $this->flushUsfmBook($state, $booksTable, $versesTable);

                $osisId = OsisBookId::normalize($matches[1]);

                if ($osisId === null) {
                    continue;
                }

                $state->pendingOsisId = $osisId;

                continue;
            }

            if ($state->pendingOsisId === null) {
                continue;
            }

            if ($state->bookDbId === null && preg_match($this->markerPattern(self::HEADING_MARKER) . '\s+(.+)$/u', $line, $matches)) {
                $state->bookDbId = $this->insertBook($booksTable, $state->pendingOsisId, trim($matches[1]));

                continue;
            }

            if (preg_match($this->markerPattern(self::CHAPTER_MARKER) . '\s+(\d+)/u', $line, $matches)) {
                $this->ensureUsfmBook($state, $booksTable);
                $this->flushUsfmVerse($state, $versesTable);
                $state->chapter = (int) $matches[1];
                $state->verse = 0;
                $state->verseText = '';
                $state->maxChapter = max($state->maxChapter, $state->chapter);

                continue;
            }

            if ($state->bookDbId === null) {
                continue;
            }

            if (preg_match($this->markerPattern(self::VERSE_MARKER) . '\s+(\d+)(?:-\d+)?\s*(.*)$/u', $line, $matches)) {
                $this->flushUsfmVerse($state, $versesTable);
                $state->verse = (int) $matches[1];
                $state->verseText = $matches[2];

                continue;
            }

            if ($state->verse < 1 || $state->chapter < 1) {
                continue;
            }

            if (str_starts_with($line, '\\')) {
                if (preg_match('/^\\\\[a-z0-9*]+\s*(.*)$/iu', $line, $matches) && trim($matches[1]) !== '') {
                    $state->verseText = trim($state->verseText . ' ' . $matches[1]);
                }

                continue;
            }

            $state->verseText = trim($state->verseText . ' ' . $line);
        }

        $this->flushUsfmBook($state, $booksTable, $versesTable);
    }

    private function flushUsfmVerse(UsfmImportState $state, string $versesTable): void
    {
        if ($state->bookDbId === null || $state->chapter < 1 || $state->verse < 1 || trim($state->verseText) === '') {
            return;
        }

        DB::table($versesTable)->insert([
            'book_id' => $state->bookDbId,
            'chapter' => $state->chapter,
            'verse' => $state->verse,
            'text' => trim($state->verseText),
            'plain_text' => $this->verseTextFormatter->toPlainText(trim($state->verseText)),
        ]);
    }

    private function flushUsfmBook(UsfmImportState $state, string $booksTable, string $versesTable): void
    {
        $this->flushUsfmVerse($state, $versesTable);

        if ($state->bookDbId !== null) {
            DB::table($booksTable)->where('id', $state->bookDbId)->update(['chapters' => $state->maxChapter]);
        }

        $state->bookDbId = null;
        $state->pendingOsisId = null;
        $state->chapter = 0;
        $state->verse = 0;
        $state->verseText = '';
        $state->maxChapter = 0;
    }

    private function ensureUsfmBook(UsfmImportState $state, string $booksTable): void
    {
        if ($state->bookDbId !== null || $state->pendingOsisId === null) {
            return;
        }

        $state->bookDbId = $this->insertBook($booksTable, $state->pendingOsisId, strtoupper($state->pendingOsisId));
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

        $items = scandir($directory);

        if ($items === false) {
            return;
        }

        foreach ($items as $item) {
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
