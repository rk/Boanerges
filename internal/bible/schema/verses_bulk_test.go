package schema_test

import (
	"testing"

	"boanerges/internal/bible/schema"
)

func TestInsertVersesBulk(t *testing.T) {
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

	rows := []schema.VerseRow{
		{BookID: bookID, Chapter: 1, Verse: 1, Text: "a", PlainText: "a"},
		{BookID: bookID, Chapter: 1, Verse: 2, Text: "b", PlainText: "b"},
	}
	if err := mgr.InsertVerses("tst", rows); err != nil {
		t.Fatal(err)
	}

	n, err := mgr.VerseCount("tst")
	if err != nil || n != 2 {
		t.Fatalf("verse count = %d, err = %v", n, err)
	}
}
