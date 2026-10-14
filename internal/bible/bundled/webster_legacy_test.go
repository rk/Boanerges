package bundled_test

import (
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/schema"
	"boanerges/internal/dictionary"
)

func TestImportWebster1828FromLegacyDB(t *testing.T) {
	legacy := filepath.Join(t.TempDir(), "legacy.sqlite")
	userDB := filepath.Join(t.TempDir(), "user.sqlite")
	seedMiniLegacyDictionaryDB(t, legacy)

	if err := bundled.ImportWebster1828FromLegacyDB(userDB, legacy); err != nil {
		t.Fatal(err)
	}
	if n := countDictionaryEntries(t, userDB); n != 2 {
		t.Fatalf("entries = %d", n)
	}

	mgr, err := schema.Open(userDB)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()
	if !dictionary.New(mgr.DB()).IsImported() {
		t.Fatal("expected imported meta")
	}
}

func TestImportWebster1828FromLegacyDB_fullReference(t *testing.T) {
	if testing.Short() {
		t.Skip("slow: full legacy nativephp dictionary copy")
	}
	legacyPath, err := bundled.FindLegacyDictionaryDB()
	if err != nil {
		t.Skip(err)
	}
	userDB := filepath.Join(t.TempDir(), "user.sqlite")
	if err := bundled.ImportWebster1828FromLegacyDB(userDB, legacyPath); err != nil {
		t.Fatal(err)
	}
	if n := countDictionaryEntries(t, userDB); n < 100_000 {
		t.Fatalf("entries = %d", n)
	}
}
