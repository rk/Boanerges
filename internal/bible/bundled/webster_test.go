package bundled_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/schema"

	_ "modernc.org/sqlite"
)

func TestImportWebster1828StaleMetaReimports(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "bible.sqlite")
	mgr, err := schema.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateAppSchema(); err != nil {
		t.Fatal(err)
	}
	_, err = mgr.DB().Exec(`INSERT INTO import_meta (key, completed_at) VALUES (?, datetime('now'))`, "webster1828_dictionary_v1")
	mgr.Close()
	if err != nil {
		t.Fatal(err)
	}
	if err := bundled.ImportWebster1828(dbPath, miniWebsterJSON(t)); err != nil {
		t.Fatal(err)
	}
	if n := countDictionaryEntries(t, dbPath); n != 2 {
		t.Fatalf("entries = %d", n)
	}
}

func TestImportWebster1828Mini(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "bible.sqlite")
	if err := bundled.ImportWebster1828(dbPath, miniWebsterJSON(t)); err != nil {
		t.Fatal(err)
	}
	if n := countDictionaryEntries(t, dbPath); n != 2 {
		t.Fatalf("entries = %d", n)
	}
}

func TestResetWebster1828(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "bible.sqlite")
	if err := bundled.ImportWebster1828(dbPath, miniWebsterJSON(t)); err != nil {
		t.Fatal(err)
	}
	if err := bundled.ResetWebster1828(dbPath); err != nil {
		t.Fatal(err)
	}
	if n := countDictionaryEntries(t, dbPath); n != 0 {
		t.Fatalf("entries after reset = %d", n)
	}
	mgr, err := schema.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()
	if bundled.WebsterInstalled(mgr.DB()) {
		t.Fatal("WebsterInstalled should be false after reset")
	}
}

func TestImportWebster1828(t *testing.T) {
	if testing.Short() {
		t.Skip("slow: full bundled dictionary.json import")
	}
	root := filepath.Join("..", "..", "..")
	jsonPath := filepath.Join(root, "extras", "websters1828", "dictionary.json")
	dbPath := filepath.Join(t.TempDir(), "bible.sqlite")

	if err := bundled.ImportWebster1828(dbPath, jsonPath); err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM dictionary_entries`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n < 1000 {
		t.Fatalf("expected many entries, got %d", n)
	}

	var completed sql.NullString
	if err := db.QueryRow(`SELECT completed_at FROM import_meta WHERE key = 'webster1828_dictionary_v1'`).Scan(&completed); err != nil || !completed.Valid {
		t.Fatal("import meta not set", err, completed)
	}
}
