<?php

namespace App\Services\Bible\Import;

use App\Enums\CatalogImportFormat;
use App\Models\Translation;
use App\Services\Bible\TranslationSchemaManager;
use Closure;
use Illuminate\Support\Facades\DB;
use InvalidArgumentException;

class UsfmFormatImporter implements TranslationFormatImporter
{
    public function __construct(
        private UsfmImporter $usfmImporter,
        private TranslationSchemaManager $schema,
    ) {}

    public function format(): CatalogImportFormat
    {
        return CatalogImportFormat::Usfm;
    }

    public function isAlreadyAvailable(Translation $translation): bool
    {
        return false;
    }

    public function prepareDownloadedSource(Translation $translation, string $downloadedPath): ?string
    {
        return $downloadedPath;
    }

    public function import(Translation $translation, ?string $sourcePath, Closure $onProgress): void
    {
        if ($sourcePath === null) {
            throw new InvalidArgumentException('Missing USFM source.');
        }

        $onProgress(70);
        $this->usfmImporter->importFromZip($translation->abbrev, $sourcePath);
    }

    public function verify(Translation $translation): void
    {
        if (! $this->hasVerse($translation->abbrev, 'gen', 1, 1) && ! $this->hasVerse($translation->abbrev, 'mat', 1, 1)) {
            throw new \RuntimeException('Verification failed: no reference verses found.');
        }
    }

    private function hasVerse(string $abbrev, string $bookId, int $chapter, int $verse): bool
    {
        $books = $this->schema->booksTable($abbrev);
        $verses = $this->schema->versesTable($abbrev);

        $book = DB::table($books)->where('osis_id', $bookId)->first();

        if ($book === null) {
            return false;
        }

        return DB::table($verses)
            ->where('book_id', $book->id)
            ->where('chapter', $chapter)
            ->where('verse', $verse)
            ->where('text', '!=', '')
            ->exists();
    }
}
