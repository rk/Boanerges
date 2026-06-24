<?php

namespace App\Services\Bible;

use App\Services\Bible\Markup\VerseTextFormatter;
use Illuminate\Support\Facades\DB;
use InvalidArgumentException;

class TranslationSchemaManager
{
    public function validateAbbrev(string $abbrev): string
    {
        $normalized = strtolower($abbrev);

        if (! preg_match('/^[a-z0-9_]+$/', $normalized)) {
            throw new InvalidArgumentException("Invalid translation abbrev: {$abbrev}");
        }

        return $normalized;
    }

    public function tableExists(string $abbrev): bool
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $books = $this->booksTable($abbrev);

        return in_array($books, $this->tableNames(), true);
    }

    public function createTables(string $abbrev): void
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $books = $this->booksTable($abbrev);
        $verses = $this->versesTable($abbrev);

        DB::statement("CREATE TABLE IF NOT EXISTS {$books} (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            osis_id TEXT NOT NULL,
            testament TEXT NOT NULL,
            chapters INTEGER NOT NULL
        )");

        DB::statement("CREATE TABLE IF NOT EXISTS {$verses} (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            book_id INTEGER NOT NULL,
            chapter INTEGER NOT NULL,
            verse INTEGER NOT NULL,
            text TEXT NOT NULL,
            plain_text TEXT NOT NULL DEFAULT '',
            FOREIGN KEY (book_id) REFERENCES {$books}(id)
        )");

        DB::statement("CREATE INDEX IF NOT EXISTS {$abbrev}_verses_lookup ON {$verses} (book_id, chapter, verse)");
    }

    public function createFtsTable(string $abbrev): void
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $verses = $this->versesTable($abbrev);
        $fts = $this->ftsTable($abbrev);

        DB::statement("CREATE VIRTUAL TABLE IF NOT EXISTS {$fts} USING fts5(
            plain_text,
            book_id UNINDEXED,
            chapter UNINDEXED,
            verse UNINDEXED,
            content='{$verses}',
            content_rowid='id'
        )");
    }

    public function ensurePlainTextColumn(string $abbrev): void
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $verses = $this->versesTable($abbrev);

        if (! in_array($verses, $this->tableNames(), true)) {
            return;
        }

        if (in_array('plain_text', $this->tableColumns($verses), true)) {
            return;
        }

        DB::statement("ALTER TABLE {$verses} ADD COLUMN plain_text TEXT NOT NULL DEFAULT ''");
    }

    public function backfillPlainText(string $abbrev, VerseTextFormatter $formatter): void
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $verses = $this->versesTable($abbrev);

        $this->ensurePlainTextColumn($abbrev);

        DB::table($verses)
            ->orderBy('id')
            ->lazyById(500)
            ->each(function ($row) use ($verses, $formatter): void {
                DB::table($verses)->where('id', $row->id)->update([
                    'plain_text' => $formatter->toPlainText((string) $row->text),
                ]);
            });
    }

    public function rebuildFtsIndex(string $abbrev, ?VerseTextFormatter $formatter = null): void
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $fts = $this->ftsTable($abbrev);

        if ($formatter !== null) {
            $this->backfillPlainText($abbrev, $formatter);
        }

        if (in_array($fts, $this->tableNames(), true)) {
            DB::statement("DROP TABLE IF EXISTS {$fts}");
        }

        $this->createFtsTable($abbrev);

        DB::statement("INSERT INTO {$fts}({$fts}) VALUES('rebuild')");
    }

    /** @return list<string> */
    private function tableColumns(string $table): array
    {
        return array_column(
            DB::select('PRAGMA table_info(' . $table . ')'),
            'name',
        );
    }

    public function dropTables(string $abbrev): void
    {
        $abbrev = $this->validateAbbrev($abbrev);

        foreach ([$this->ftsTable($abbrev), $this->versesTable($abbrev), $this->booksTable($abbrev)] as $table) {
            if (in_array($table, $this->tableNames(), true)) {
                DB::statement("DROP TABLE IF EXISTS {$table}");
            }
        }
    }

    public function hasBook(string $abbrev, int $bookId): bool
    {
        $abbrev = $this->validateAbbrev($abbrev);
        $books = $this->booksTable($abbrev);

        if (! in_array($books, $this->tableNames(), true)) {
            return false;
        }

        return DB::table($books)->where('id', $bookId)->exists();
    }

    public function booksTable(string $abbrev): string
    {
        return $this->validateAbbrev($abbrev) . '_books';
    }

    public function versesTable(string $abbrev): string
    {
        return $this->validateAbbrev($abbrev) . '_verses';
    }

    public function ftsTable(string $abbrev): string
    {
        return $this->validateAbbrev($abbrev) . '_verses_fts';
    }

    /** @return list<string> */
    private function tableNames(): array
    {
        return array_column(
            DB::select("SELECT name FROM sqlite_master WHERE type='table'"),
            'name',
        );
    }
}
