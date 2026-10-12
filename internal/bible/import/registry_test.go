package bibleimport_test

import (
	"testing"

	"boanerges/internal/bible/catalog"
	bibleimport "boanerges/internal/bible/import"
	"boanerges/internal/bible/import/swordfmt"
	"boanerges/internal/bible/import/usfm"
	"boanerges/internal/bible/schema"
)

func TestEveryCatalogImportFormatResolves(t *testing.T) {
	mgr, err := schema.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	reg := bibleimport.NewRegistry(
		swordfmt.NewImporter(mgr, nil, nil),
		usfm.NewImporter(mgr, nil),
	)

	for _, format := range []catalog.ImportFormat{catalog.FormatSword, catalog.FormatUsfm} {
		imp, err := reg.Format(format)
		if err != nil {
			t.Fatal(err)
		}
		if imp.Format() != format {
			t.Fatalf("format %s resolved to %s", format, imp.Format())
		}
	}

	if _, err := reg.Format(catalog.FormatAccordance); err == nil {
		t.Fatal("expected Accordance format to be unsupported")
	}
}
