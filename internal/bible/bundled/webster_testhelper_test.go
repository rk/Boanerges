package bundled_test

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"boanerges/internal/bible/schema"

	_ "modernc.org/sqlite"
)

func seedMiniLegacyDictionaryDB(t *testing.T, dbPath string) {
	t.Helper()
	mgr, err := schema.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()
	if err := mgr.CreateAppSchema(); err != nil {
		t.Fatal(err)
	}
	db := mgr.DB()
	if _, err := db.Exec(`INSERT INTO dictionary_entries (word, word_key, part_of_speech, definitions) VALUES
		('Love', 'love', 'n.', '["Strong affection."]'),
		('Faith', 'faith', 'n.', '["Trust."]')`); err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec(`INSERT INTO import_meta (key, completed_at) VALUES (?, ?)`,
		"webster1828_dictionary_v1", time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatal(err)
	}
}

func countDictionaryEntries(t *testing.T, dbPath string) int {
	t.Helper()
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM dictionary_entries`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	return n
}

func miniWebsterJSON(t *testing.T) string {
	t.Helper()
	return filepath.Join("testdata", "webster-mini.json")
}
