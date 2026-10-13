package read_test

import (
	"testing"

	"boanerges/internal/bible/read"
	"boanerges/internal/bible/schema"
)

func TestListBooksSkipsEmptyCoverage(t *testing.T) {
	mgr, err := schema.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	if err := mgr.CreateTables("bk"); err != nil {
		t.Fatal(err)
	}
	res, err := mgr.DB().Exec(`INSERT INTO bk_books (name, osis_id, testament, chapters) VALUES ('Gen','gen','ot',1)`)
	if err != nil {
		t.Fatal(err)
	}
	genID, _ := res.LastInsertId()
	if _, err := mgr.DB().Exec(`INSERT INTO bk_books (name, osis_id, testament, chapters) VALUES ('Empty','empty','ot',1)`); err != nil {
		t.Fatal(err)
	}
	if err := mgr.InsertVerses("bk", []schema.VerseRow{
		{BookID: genID, Chapter: 1, Verse: 1, Text: "start", PlainText: "start"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := mgr.IndexBookCoverage("bk"); err != nil {
		t.Fatal(err)
	}

	books, err := read.ListBooks(mgr.DB(), mgr, "bk")
	if err != nil {
		t.Fatal(err)
	}
	if len(books) != 1 || books[0].ID != "gen" {
		t.Fatalf("books = %+v", books)
	}
}
