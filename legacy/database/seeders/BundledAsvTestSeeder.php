<?php

namespace Database\Seeders;

use App\Enums\TranslationInstallStatus;
use App\Enums\TranslationInstallStep;
use App\Enums\VerseMarkupFormat;
use App\Models\Translation;
use App\Services\Bible\Markup\VerseMarkupConverterFactory;
use App\Services\Bible\TranslationSchemaManager;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

/**
 * Minimal ASV fixture for feature tests (fast; no SWORD import).
 *
 * Phase 2 (future): replace body with loading tests/fixtures/asv-bundled.sql exported
 * from a real import so tests use authentic verse text without parse cost.
 */
class BundledAsvTestSeeder extends Seeder
{
    public function run(): void
    {
        if (Translation::query()
            ->where('abbrev', 'asv')
            ->where('install_status', TranslationInstallStatus::Ready)
            ->exists()) {
            return;
        }

        $schema = app(TranslationSchemaManager::class);
        $formatter = VerseMarkupConverterFactory::defaultFormatter();

        Translation::query()->create([
            'abbrev' => 'asv',
            'name' => 'American Standard Version',
            'format' => VerseMarkupFormat::Osis->value,
            'install_status' => TranslationInstallStatus::Ready,
            'install_step' => TranslationInstallStep::Ready,
            'bundled' => true,
        ]);

        $schema->createTables('asv');

        DB::table($schema->booksTable('asv'))->insert([
            'id' => 1,
            'name' => 'Genesis',
            'osis_id' => 'gen',
            'testament' => 'OT',
            'chapters' => 50,
        ]);

        $verses = [];

        for ($verse = 1; $verse <= 31; $verse++) {
            $text = $verse === 1
                ? 'In the beginning God created the heaven and the earth.'
                : "Genesis 1:{$verse} placeholder text.";

            $verses[] = [
                'book_id' => 1,
                'chapter' => 1,
                'verse' => $verse,
                'text' => $text,
                'plain_text' => $formatter->toPlainText($text),
            ];
        }

        DB::table($schema->versesTable('asv'))->insert($verses);

        $schema->createFtsTable('asv');
        DB::statement('INSERT INTO ' . $schema->ftsTable('asv') . '(' . $schema->ftsTable('asv') . ") VALUES('rebuild')");
    }
}
