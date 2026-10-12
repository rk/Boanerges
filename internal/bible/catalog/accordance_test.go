package catalog_test

import (
	"path/filepath"
	"testing"

	"boanerges/internal/bible/catalog"
)

func TestCatalogHasNoAccordanceEntries(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	cat, err := catalog.Load(filepath.Join(root, "extras", "translations.json"))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range cat.All() {
		if entry.ImportAs == catalog.FormatAccordance {
			t.Fatalf("catalog entry %q uses Accordance import (not ported)", entry.Short)
		}
	}
}
