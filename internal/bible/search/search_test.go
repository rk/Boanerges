package search_test

import (
	"testing"

	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/search"
	"boanerges/internal/bible/translations"
)

func seedSearchDB(t *testing.T, abbrev string) (*schema.Manager, *translations.Store) {
	t.Helper()
	mgr, err := schema.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if !schema.SupportsFTS5(mgr.DB()) {
		t.Skip("no fts5")
	}
	if err := mgr.CreateAppSchema(); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateTables(abbrev); err != nil {
		t.Fatal(err)
	}
	res, err := mgr.DB().Exec(`INSERT INTO `+abbrev+`_books (name, osis_id, testament, chapters) VALUES ('Gen','gen','ot',1)`)
	if err != nil {
		t.Fatal(err)
	}
	bookID, _ := res.LastInsertId()
	if err := mgr.InsertVerses(abbrev, []schema.VerseRow{
		{BookID: bookID, Chapter: 1, Verse: 1, Text: "In the beginning", PlainText: "In the beginning"},
		{BookID: bookID, Chapter: 1, Verse: 2, Text: "And the earth", PlainText: "And the earth"},
		{BookID: bookID, Chapter: 1, Verse: 3, Text: "beginning again", PlainText: "beginning again"},
	}); err != nil {
		t.Fatal(err)
	}
	if err := mgr.RebuildFtsIndex(abbrev, markup.DefaultFormatter()); err != nil {
		t.Fatal(err)
	}
	store := translations.NewStore(mgr)
	if _, err := mgr.DB().Exec(`
		INSERT INTO translations (abbrev, name, install_status, install_step, bundled, created_at, updated_at)
		VALUES (?, 'Test', 'ready', 'ready', 0, datetime('now'), datetime('now'))`, abbrev); err != nil {
		t.Fatal(err)
	}

	return mgr, store
}

func TestSearchPaginationHasMore(t *testing.T) {
	mgr, store := seedSearchDB(t, "pag")
	defer mgr.Close()

	abbrev := "pag"
	resp, err := search.Search(mgr.DB(), mgr, store, "beginning", &abbrev, 1, 0)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Total < 2 {
		t.Fatalf("total = %d", resp.Total)
	}
	if !resp.HasMore || len(resp.Results) != 1 {
		t.Fatalf("page1 = %+v", resp)
	}

	resp2, err := search.Search(mgr.DB(), mgr, store, "beginning", &abbrev, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if resp2.HasMore && len(resp2.Results) == 0 {
		t.Fatalf("page2 empty but hasMore: %+v", resp2)
	}
}
