package install

import (
	"testing"

	"boanerges/internal/bible/schema"
)

func TestImportLooksComplete(t *testing.T) {
	mgr, err := schema.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	if err := mgr.CreateTables("tst"); err != nil {
		t.Fatal(err)
	}
	res, err := mgr.DB().Exec(`INSERT INTO tst_books (name, osis_id, testament, chapters) VALUES ('Gen','gen','ot',1)`)
	if err != nil {
		t.Fatal(err)
	}
	bookID, _ := res.LastInsertId()

	rows := make([]schema.VerseRow, 0, 1000)
	for v := 1; v <= 1000; v++ {
		rows = append(rows, schema.VerseRow{
			BookID: bookID, Chapter: 1, Verse: v, Text: "x", PlainText: "x",
		})
	}
	if err := mgr.InsertVerses("tst", rows); err != nil {
		t.Fatal(err)
	}

	p := &Pipeline{schema: mgr}
	if !p.importLooksComplete("tst") {
		t.Fatal("expected import complete with 1000 verses and gen 1:1")
	}
}

func TestRebuildFtsFromVersesTable(t *testing.T) {
	mgr, err := schema.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()
	if !schema.SupportsFTS5(mgr.DB()) {
		t.Skip("no fts5")
	}

	if err := mgr.CreateTables("tst"); err != nil {
		t.Fatal(err)
	}
	res, err := mgr.DB().Exec(`INSERT INTO tst_books (name, osis_id, testament, chapters) VALUES ('Gen','gen','ot',1)`)
	if err != nil {
		t.Fatal(err)
	}
	bookID, _ := res.LastInsertId()
	if err := mgr.InsertVerses("tst", []schema.VerseRow{
		{BookID: bookID, Chapter: 1, Verse: 1, Text: "In the beginning", PlainText: "In the beginning"},
	}); err != nil {
		t.Fatal(err)
	}

	if err := mgr.RebuildFtsIndex("tst", nil); err != nil {
		t.Fatal(err)
	}
	ok, err := mgr.FtsPopulated("tst")
	if err != nil || !ok {
		t.Fatalf("fts populated = %v err = %v", ok, err)
	}
}
