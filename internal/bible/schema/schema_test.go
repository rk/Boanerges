package schema

import (
	"database/sql"
	"path/filepath"
	"testing"

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

func TestValidateAbbrev(t *testing.T) {
	mgr, _ := Open(filepath.Join(t.TempDir(), "x.db"))
	defer mgr.Close()

	if _, err := mgr.ValidateAbbrev("Bad-Id"); err == nil {
		t.Fatal("expected invalid abbrev error")
	}
}
