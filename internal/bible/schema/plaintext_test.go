package schema_test

import (
	"testing"

	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
)

func TestBackfillPlainTextBulkUpdate(t *testing.T) {
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
	if err := mgr.InsertVerses("tst", []schema.VerseRow{
		{BookID: bookID, Chapter: 1, Verse: 1, Text: "<i>Hello</i>", PlainText: ""},
		{BookID: bookID, Chapter: 1, Verse: 2, Text: "World", PlainText: "World"},
	}); err != nil {
		t.Fatal(err)
	}

	if err := mgr.BackfillPlainTextIfEmpty("tst", markup.DefaultFormatter()); err != nil {
		t.Fatal(err)
	}

	var plain string
	if err := mgr.DB().QueryRow(`SELECT plain_text FROM tst_verses WHERE verse = 1`).Scan(&plain); err != nil {
		t.Fatal(err)
	}
	if plain == "" {
		t.Fatal("expected plain_text backfill for verse 1")
	}
}
