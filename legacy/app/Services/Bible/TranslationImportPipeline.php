<?php

namespace App\Services\Bible;

use App\Enums\CatalogImportFormat;
use App\Enums\TranslationInstallStatus;
use App\Enums\TranslationInstallStep;
use App\Events\FtsIndexProgress;
use App\Models\Translation;
use App\Services\Bible\Import\TranslationImporterRegistry;
use App\Services\Bible\Markup\VerseTextFormatter;
use Illuminate\Support\Facades\Http;
use Illuminate\Support\Facades\Storage;

class TranslationImportPipeline
{
    private ?string $downloadedPath = null;

    public function __construct(
        private TranslationCatalog $catalog,
        private TranslationSchemaManager $schema,
        private TranslationImporterRegistry $importers,
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
        $importer = $this->importerFor($translation);

        if ($importer->isAlreadyAvailable($translation)) {
            $translation->updateProgress(TranslationInstallStatus::Downloading, TranslationInstallStep::SourceReady, 10);
            $this->syncMetadata($translation);

            return;
        }

        $translation->updateProgress(TranslationInstallStatus::Downloading, TranslationInstallStep::Downloading, 0);

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

        $preparedPath = $importer->prepareDownloadedSource($translation, $zipPath);
        $this->downloadedPath = $preparedPath;

        $translation->updateProgress(TranslationInstallStatus::Downloading, TranslationInstallStep::Downloaded, 10);

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
        $translation->updateProgress(TranslationInstallStatus::CreatingSchema, TranslationInstallStep::CreatingSchema, 20);
    }

    public function importVerses(Translation $translation): void
    {
        $importer = $this->importerFor($translation);

        $importer->import(
            $translation,
            $this->downloadedPath,
            function (float|int $percent) use ($translation): void {
                $translation->updateProgress(
                    TranslationInstallStatus::Importing,
                    TranslationInstallStep::Importing,
                    (int) round($percent),
                );
            },
        );
    }

    public function verify(Translation $translation): void
    {
        $this->importerFor($translation)->verify($translation);

        $translation->updateProgress(TranslationInstallStatus::Verifying, TranslationInstallStep::Verifying, 75);
    }

    public function buildFtsIndex(Translation $translation): void
    {
        $translation->updateProgress(TranslationInstallStatus::Indexing, TranslationInstallStep::Indexing, 85);
        $this->schema->rebuildFtsIndex($translation->abbrev, $this->verseTextFormatter);

        $translation->updateProgress(TranslationInstallStatus::Indexing, TranslationInstallStep::Indexing, 95);

        event(new FtsIndexProgress(
            abbrev: $translation->abbrev,
            step: TranslationInstallStep::Indexed,
            percent: 95,
        ));
    }

    public function markReady(Translation $translation): void
    {
        $translation->update([
            'install_status' => TranslationInstallStatus::Ready,
            'install_step' => TranslationInstallStep::Ready,
            'install_error' => null,
        ]);

        event(new \App\Events\TranslationInstallProgress(
            abbrev: $translation->abbrev,
            step: TranslationInstallStep::Ready,
            percent: 100,
        ));
    }

    private function importerFor(Translation $translation): \App\Services\Bible\Import\TranslationFormatImporter
    {
        return $this->importers->get($this->importAs($translation));
    }

    private function importAs(Translation $translation): CatalogImportFormat
    {
        return $this->catalog->find($translation->abbrev)->importAs;
    }
}
