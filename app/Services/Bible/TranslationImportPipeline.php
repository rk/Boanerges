<?php

namespace App\Services\Bible;

use App\Enums\TranslationInstallStatus;
use App\Events\FtsIndexProgress;
use App\Models\Translation;
use App\Services\Bible\Import\AccordanceImporter;
use App\Services\Bible\Import\SwordImporter;
use App\Services\Bible\Import\UsfmImporter;
use App\Services\Bible\Markup\VerseTextFormatter;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Storage;
use ZipArchive;

class TranslationImportPipeline
{
    private ?string $downloadedPath = null;

    public function __construct(
        private BibleModuleManager $modules,
        private TranslationCatalog $catalog,
        private TranslationSchemaManager $schema,
        private SwordImporter $swordImporter,
        private UsfmImporter $usfmImporter,
        private AccordanceImporter $accordanceImporter,
        private VerseTextFormatter $verseTextFormatter,
        private TranslationMetadataSync $metadataSync,
    ) {}

    public function run(Translation $translation): void
    {
        $this->downloadedPath = null;

        try {
            $this->download($translation);
            $this->createSchema($translation);
            $this->importVerses($translation);
            $this->verify($translation);
            $this->buildFtsIndex($translation);
            $this->markReady($translation);
        } catch (\Throwable $exception) {
            $translation->markFailed($exception->getMessage());

            throw $exception;
        } finally {
            if ($this->downloadedPath !== null && is_file($this->downloadedPath)) {
                @unlink($this->downloadedPath);
            }
        }
    }

    public function download(Translation $translation): void
    {
        $importAs = $this->importAs($translation);

        if ($importAs === 'sword' && $this->modules->isModuleInstalled($translation->abbrev)) {
            $translation->updateProgress(TranslationInstallStatus::Downloading, 'source_ready', 10);
            $this->syncMetadata($translation);

            return;
        }

        $translation->updateProgress(TranslationInstallStatus::Downloading, 'downloading', 0);

        $entry = $this->catalog->find($translation->abbrev);
        $zipPath = Storage::disk('local')->path('tmp/' . $entry->short . '.module.zip');

        if (! is_dir(dirname($zipPath))) {
            mkdir(dirname($zipPath), 0755, true);
        }

        $response = Http::timeout(120)->sink($zipPath)->get($entry->url);

        if (! $response->successful() || ! is_file($zipPath) || filesize($zipPath) === 0) {
            @unlink($zipPath);
            throw new \RuntimeException("Failed to download {$entry->short}.");
        }

        if ($importAs === 'sword') {
            $zip = new ZipArchive();

            if ($zip->open($zipPath) !== true) {
                @unlink($zipPath);
                throw new \RuntimeException("Failed to open archive for {$entry->short}.");
            }

            $zip->extractTo($this->modules->localRoot());
            $zip->close();
            @unlink($zipPath);
            $this->modules->clearCache();
        } else {
            $this->downloadedPath = $zipPath;
        }

        $translation->updateProgress(TranslationInstallStatus::Downloading, 'downloaded', 10);

        $this->syncMetadata($translation);
    }

    private function syncMetadata(Translation $translation): void
    {
        $entry = $this->catalog->find($translation->abbrev);
        $this->metadataSync->applyFromCatalog($translation, $entry);
    }

    public function createSchema(Translation $translation): void
    {
        $this->schema->dropTables($translation->abbrev);
        $this->schema->createTables($translation->abbrev);
        $translation->updateProgress(TranslationInstallStatus::CreatingSchema, 'creating_schema', 20);
    }

    public function importVerses(Translation $translation): void
    {
        switch ($this->importAs($translation)) {
            case 'usfm':
                $translation->updateProgress(TranslationInstallStatus::Importing, 'importing', 70);
                $this->usfmImporter->importFromZip(
                    $translation->abbrev,
                    $this->downloadedPath ?? throw new \RuntimeException('Missing USFM source.'),
                );
                break;
            case 'accordance':
                $translation->updateProgress(TranslationInstallStatus::Importing, 'importing', 70);
                $this->accordanceImporter->importFromFile(
                    $translation->abbrev,
                    $this->downloadedPath ?? throw new \RuntimeException('Missing Accordance source.'),
                );
                break;
            default:
                $this->swordImporter->progressConfigure(
                    20,
                    50,
                    static function (float $percent) use ($translation) {
                        $translation->updateProgress(TranslationInstallStatus::Importing, 'importing', round($percent));
                    },
                );
                $this->swordImporter->import($translation->abbrev, $translation->abbrev);
                break;
        }
    }

    public function verify(Translation $translation): void
    {
        if ($this->importAs($translation) === 'sword') {
            $this->swordImporter->verify($translation->abbrev);
        } elseif (! $this->hasVerse($translation->abbrev, 'gen', 1, 1) && ! $this->hasVerse($translation->abbrev, 'mat', 1, 1)) {
            throw new \RuntimeException('Verification failed: no reference verses found.');
        }

        $translation->updateProgress(TranslationInstallStatus::Verifying, 'verifying', 75);
    }

    public function buildFtsIndex(Translation $translation): void
    {
        $translation->updateProgress(TranslationInstallStatus::Indexing, 'indexing', 85);
        $this->schema->rebuildFtsIndex($translation->abbrev, $this->verseTextFormatter);

        $translation->updateProgress(TranslationInstallStatus::Indexing, 'indexing', 95);

        event(new FtsIndexProgress(
            abbrev: $translation->abbrev,
            step: 'indexed',
            percent: 95,
        ));
    }

    public function markReady(Translation $translation): void
    {
        $translation->update([
            'install_status' => TranslationInstallStatus::Ready,
            'install_step' => 'ready',
            'install_error' => null,
        ]);

        event(new \App\Events\TranslationInstallProgress(
            abbrev: $translation->abbrev,
            step: 'ready',
            percent: 100,
        ));
    }

    private function importAs(Translation $translation): string
    {
        return $this->catalog->find($translation->abbrev)->importAs;
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
