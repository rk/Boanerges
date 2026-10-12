package bundled_test

import (
	"database/sql"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"

	_ "modernc.org/sqlite"
)

func TestImportASVFromBundledSQLite(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	bundledPath := filepath.Join(root, "extras", "bundled", "asv.sqlite")
	confPath := filepath.Join(root, "extras", "sword", "mods.d", "asv.conf")

	userDB := filepath.Join(t.TempDir(), "app.db")
	if err := bundled.ImportASV(userDB, bundledPath, confPath); err != nil {
		t.Skip("bundled ASV sqlite not built:", err)
	}

	db, err := sql.Open("sqlite", userDB)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	var verses int
	if err := db.QueryRow(`SELECT COUNT(*) FROM asv_verses`).Scan(&verses); err != nil {
		t.Fatal(err)
	}
	if verses < 30000 {
		t.Fatalf("expected full bible verse count, got %d", verses)
	}

	var abbrev string
	if err := db.QueryRow(`SELECT abbrev FROM translations WHERE abbrev='asv'`).Scan(&abbrev); err != nil {
		t.Fatal(err)
	}

	// idempotent
	if err := bundled.ImportASV(userDB, bundledPath, confPath); err != nil {
		t.Fatal(err)
	}
	var verses2 int
	if err := db.QueryRow(`SELECT COUNT(*) FROM asv_verses`).Scan(&verses2); err != nil {
		t.Fatal(err)
	}
	if verses2 != verses {
		t.Fatalf("re-import changed count: %d -> %d", verses, verses2)
	}
}
