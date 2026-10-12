package bundled_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"

	_ "modernc.org/sqlite"
)

func TestImportWebster1828(t *testing.T) {
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
