package catalog_test

import (
	"path/filepath"
	"testing"

	"boanerges/internal/bible/catalog"
)

func TestCatalogLoadsEnglishTranslations(t *testing.T) {
	path := filepath.Join("..", "..", "..", "extras", "translations.json")
	cat, err := catalog.Load(path)
	if err != nil {
		t.Fatal(err)
	}

	entries := cat.All()
	if len(entries) < 50 {
		t.Fatalf("expected >50 catalog entries, got %d", len(entries))
	}

	asv, err := cat.Find("ASV")
	if err != nil {
		t.Fatal(err)
	}
	if asv.Name == "" {
		t.Fatal("ASV entry missing name")
	}
}
