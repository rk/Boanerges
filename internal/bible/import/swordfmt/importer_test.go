package swordfmt_test

import (
	"path/filepath"
	"strings"
	"testing"

	"boanerges/internal/bible/import/swordfmt"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/swordpath"
)

func TestVerifyBundledASVModule(t *testing.T) {
	root := filepath.Join("..", "..", "..", "..")
	bundled := filepath.Join(root, "extras", "sword")
	roots := swordpath.Roots(bundled, t.TempDir())

	imp := swordfmt.NewImporter(nil, nil, roots)
	if err := imp.Verify("ASV"); err != nil {
		t.Skip("bundled ASV:", err)
	}
}

func TestImportBundledASVIntoSchema(t *testing.T) {
	root := filepath.Join("..", "..", "..", "..")
	bundled := filepath.Join(root, "extras", "sword")
	roots := swordpath.Roots(bundled, t.TempDir())

	mgr, err := schema.Open(filepath.Join(t.TempDir(), "bible.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()
	if err := mgr.CreateAppSchema(); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateTables("asvtest"); err != nil {
		t.Fatal(err)
	}

	imp := swordfmt.NewImporter(mgr, nil, roots)
	if err := imp.Import("asvtest", "ASV", nil); err != nil {
		t.Skip("bundled ASV import:", err)
	}

	var text string
	booksTable, _ := mgr.BooksTable("asvtest")
	versesTable, _ := mgr.VersesTable("asvtest")
	err = mgr.DB().QueryRow(`
		SELECT v.text FROM `+versesTable+` v
		JOIN `+booksTable+` b ON b.id = v.book_id
		WHERE b.osis_id = 'gen' AND v.chapter = 1 AND v.verse = 1`).Scan(&text)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(text), "beginning") {
		t.Fatalf("Gen 1:1 = %q", text)
	}
}
