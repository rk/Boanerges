package schema

import (
	"testing"
)

func TestIndexBookCoverage(t *testing.T) {
	mgr, err := Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	if err := mgr.CreateTables("cov"); err != nil {
		t.Fatal(err)
	}

	res, err := mgr.DB().Exec(`INSERT INTO cov_books (name, osis_id, testament, chapters) VALUES ('Gen','gen','ot',50)`)
	if err != nil {
		t.Fatal(err)
	}
	genID, _ := res.LastInsertId()
	res, err = mgr.DB().Exec(`INSERT INTO cov_books (name, osis_id, testament, chapters) VALUES ('Empty','empty','ot',10)`)
	if err != nil {
		t.Fatal(err)
	}
	emptyID, _ := res.LastInsertId()

	if err := mgr.InsertVerses("cov", []VerseRow{
		{BookID: genID, Chapter: 1, Verse: 1, Text: "In the beginning", PlainText: "In the beginning"},
		{BookID: genID, Chapter: 2, Verse: 1, Text: "More", PlainText: "More"},
		{BookID: emptyID, Chapter: 2, Verse: 1, Text: "skip", PlainText: "skip"},
	}); err != nil {
		t.Fatal(err)
	}

	if err := mgr.IndexBookCoverage("cov"); err != nil {
		t.Fatal(err)
	}

	var genContent, emptyContent, genChapters int
	if err := mgr.DB().QueryRow(`SELECT has_content, chapters FROM cov_books WHERE osis_id = 'gen'`).Scan(&genContent, &genChapters); err != nil {
		t.Fatal(err)
	}
	if genContent != 1 || genChapters != 2 {
		t.Fatalf("gen = has_content %d chapters %d", genContent, genChapters)
	}
	if err := mgr.DB().QueryRow(`SELECT has_content FROM cov_books WHERE osis_id = 'empty'`).Scan(&emptyContent); err != nil {
		t.Fatal(err)
	}
	if emptyContent != 0 {
		t.Fatalf("empty book has_content = %d", emptyContent)
	}
}
