package bundled_test

import (
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/dictionary"
	"boanerges/internal/bible/schema"
)

func TestImportWebster1828FromLegacyDB(t *testing.T) {
	if testing.Short() {
		t.Skip("slow: copies full Webster table from legacy sqlite")
	}
	root := filepath.Join("..", "..", "..")
	legacy := filepath.Join(root, "legacy", "database", "nativephp.sqlite")
	if _, err := bundled.FindLegacyDictionaryDB(); err != nil {
		t.Skip("reference legacy database not present")
	}

	userDB := filepath.Join(t.TempDir(), "bible.sqlite")
	if err := bundled.ImportWebster1828FromLegacyDB(userDB, legacy); err != nil {
		t.Fatal(err)
	}

	mgr, err := schema.Open(userDB)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	if !dictionary.New(mgr.DB()).IsImported() {
		t.Fatal("expected imported meta")
	}
	var n int
	if err := mgr.DB().QueryRow(`SELECT COUNT(*) FROM dictionary_entries`).Scan(&n); err != nil || n < 1000 {
		t.Fatalf("entries = %d", n)
	}
}
