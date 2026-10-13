package schema

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"boanerges/internal/bible/markup"

	_ "modernc.org/sqlite"
)

func TestSchemaAndFTS5(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	mgr, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mgr.Close() })

	db, err := sql.Open("sqlite", path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })

	if !SupportsFTS5(db) {
		t.Fatal("expected modernc sqlite FTS5 support")
	}

	if err := mgr.CreateTables("asv"); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateFtsTable("asv"); err != nil {
		t.Fatal(err)
	}

	books, err := mgr.BooksTable("asv")
	if err != nil || books != "asv_books" {
		t.Fatalf("books table = %q err=%v", books, err)
	}

	if err := mgr.DropTables("asv"); err != nil {
		t.Fatal(err)
	}
}

func TestRebuildFtsSkipsPlainTextBackfill(t *testing.T) {
	mgr, err := Open(filepath.Join(t.TempDir(), "fts.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	if err := mgr.CreateTables("bsb"); err != nil {
		t.Fatal(err)
	}
	_, err = mgr.DB().Exec(`INSERT INTO bsb_books (osis_id, name, testament, chapters) VALUES ('gen', 'Genesis', 'ot', 1)`)
	if err != nil {
		t.Fatal(err)
	}
	var bookID int64
	if err := mgr.DB().QueryRow(`SELECT id FROM bsb_books`).Scan(&bookID); err != nil {
		t.Fatal(err)
	}
	_, err = mgr.DB().Exec(
		`INSERT INTO bsb_verses (book_id, chapter, verse, text, plain_text) VALUES (?, 1, 1, '<em>Hi</em>', 'Hi')`,
		bookID,
	)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now()
	if err := mgr.RebuildFtsIndex("bsb", markup.DefaultFormatter()); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) > 2*time.Second {
		t.Fatalf("rebuild took too long with plain_text already set: %s", time.Since(start))
	}
}

func TestValidateAbbrev(t *testing.T) {
	mgr, _ := Open(filepath.Join(t.TempDir(), "x.db"))
	defer mgr.Close()

	if _, err := mgr.ValidateAbbrev("Bad-Id"); err == nil {
		t.Fatal("expected invalid abbrev error")
	}
}
