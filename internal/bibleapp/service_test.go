package bibleapp_test

import (
	"path/filepath"
	"strings"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bibleapp"
	"boanerges/internal/dictionary"
)

func seedASV(t *testing.T, dbPath string) {
	t.Helper()
	root := filepath.Join("..", "..")
	if err := bundled.ImportASV(
		dbPath,
		filepath.Join(root, "extras", "bundled", "asv.sqlite"),
		filepath.Join(root, "extras", "sword", "mods.d", "asv.conf"),
	); err != nil {
		t.Skip("bundled ASV:", err)
	}
}

func TestBooksGenesis(t *testing.T) {
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)

	svc := bibleapp.NewBibleService(db, bootstrap.NewService())
	resp, err := svc.Books("asv")
	if err != nil {
		t.Fatal(err)
	}
	var gen *struct {
		ID           string
		Chapters     int
		FirstChapter int
		LastChapter  int
	}
	for _, b := range resp.Books {
		if b.ID == "gen" {
			gen = &struct {
				ID           string
				Chapters     int
				FirstChapter int
				LastChapter  int
			}{b.ID, b.Chapters, b.FirstChapter, b.LastChapter}
			break
		}
	}
	if gen == nil || gen.Chapters != 50 || gen.FirstChapter != 1 || gen.LastChapter != 50 {
		t.Fatalf("genesis = %+v", gen)
	}
}

func TestGenesisChapterOne(t *testing.T) {
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)

	svc := bibleapp.NewBibleService(db, bootstrap.NewService())
	resp, err := svc.Chapter("asv", "gen", 1)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Chapter.Book != "Genesis" || resp.Chapter.Chapter != 1 || len(resp.Chapter.Verses) != 31 {
		t.Fatalf("chapter = %+v", resp.Chapter)
	}
	if !strings.Contains(strings.ToLower(resp.Chapter.Verses[0].Text), "beginning") {
		t.Fatal("missing beginning")
	}
}

func TestSearchBeginning(t *testing.T) {
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)

	svc := bibleapp.NewBibleService(db, bootstrap.NewService())
	resp, err := svc.Search("beginning", "asv", 50, 0)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, r := range resp.Results {
		if r.BookID == "gen" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected gen hit")
	}
}

func TestCrossReferencesAfterImport(t *testing.T) {
	root := filepath.Join("..", "..")
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)

	crossPath := filepath.Join(root, "extras", "cross-references", "cross_references.txt")
	if err := bundled.ImportCrossReferences(db, crossPath); err != nil {
		t.Skip(err)
	}

	svc := bibleapp.NewBibleService(db, bootstrap.NewService())
	resp, err := svc.CrossReferences("gen", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(resp.References) == 0 {
		t.Fatal("expected cross refs")
	}
	if resp.References[0].BookName == "" {
		t.Fatal("missing book name")
	}
}

func TestBooksUnknownTranslation(t *testing.T) {
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)
	svc := bibleapp.NewBibleService(db, bootstrap.NewService())
	_, err := svc.Books("no-such-module")
	if err == nil {
		t.Fatal("expected error for unknown translation")
	}
}

func TestChapterOutOfRange(t *testing.T) {
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)
	svc := bibleapp.NewBibleService(db, bootstrap.NewService())
	_, err := svc.Chapter("asv", "gen", 999)
	if err == nil {
		t.Fatal("expected error for invalid chapter")
	}
}

func TestDictionaryGrace(t *testing.T) {
	root := filepath.Join("..", "..")
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	seedASV(t, db)

	if err := bundled.ImportWebster1828(db, filepath.Join(root, "extras", "websters1828", "dictionary.json")); err != nil {
		t.Skip(err)
	}

	mgr, err := schema.Open(db)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	dict := dictionary.New(mgr.DB())
	if !dict.IsImported() {
		t.Fatal("dictionary not imported")
	}
	result := dict.Lookup("grace")
	if !result.Found || len(result.Variants) < 2 {
		t.Fatalf("grace lookup = %+v", result)
	}
}
