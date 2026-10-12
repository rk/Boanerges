package usfm_test

import (
	"os"
	"strings"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/import/usfm"
	"boanerges/internal/bible/schema"
)

func TestParsesUSFMBookMarkers(t *testing.T) {
	mgr := openSchema(t, "usfmtest")
	defer mgr.Close()

	path := filepath.Join(t.TempDir(), "sample.usfm")
	if err := os.WriteFile(path, []byte(`\id MAT
\h Matthew
\c 1
\p
\v 1 A genealogy of Jesus Christ, a descendant of David and Abraham.
\v 2 Abraham was the father of Isaac,
Isaac of Jacob,
\q2 Jacob of Judah and his brothers,
`), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := usfm.NewImporter(mgr, nil).ImportFromFile("usfmtest", path, true); err != nil {
		t.Fatal(err)
	}

	var osis, name string
	if err := mgr.DB().QueryRow(`SELECT osis_id, name FROM usfmtest_books`).Scan(&osis, &name); err != nil {
		t.Fatal(err)
	}
	if osis != "mat" || name != "Matthew" {
		t.Fatalf("book = %s %s", osis, name)
	}

	rows, err := mgr.DB().Query(`SELECT text, plain_text FROM usfmtest_verses ORDER BY verse`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var texts []string
	var plains []string
	for rows.Next() {
		var text, plain string
		if err := rows.Scan(&text, &plain); err != nil {
			t.Fatal(err)
		}
		texts = append(texts, text)
		plains = append(plains, plain)
	}
	if len(texts) != 2 {
		t.Fatalf("verse count = %d", len(texts))
	}
	if !contains(texts[0], "genealogy") {
		t.Fatalf("v1 text = %q", texts[0])
	}
	if !contains(texts[1], "Jacob of Judah") {
		t.Fatalf("v2 text = %q", texts[1])
	}
	if contains(plains[1], `\q2`) {
		t.Fatalf("plain should strip markers: %q", plains[1])
	}
}

func TestSkipsNonScripturePeripherals(t *testing.T) {
	mgr := openSchema(t, "usfmperiph")
	defer mgr.Close()

	path := filepath.Join(t.TempDir(), "frt.usfm")
	_ = os.WriteFile(path, []byte(`\id FRT
\mt1 Open English Bible
\c 1
\v 1 Not scripture
`), 0o644)

	if err := usfm.NewImporter(mgr, nil).ImportFromFile("usfmperiph", path, true); err != nil {
		t.Fatal(err)
	}

	var n int
	if err := mgr.DB().QueryRow(`SELECT COUNT(*) FROM usfmperiph_books`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("expected 0 books, got %d", n)
	}
}

func openSchema(t *testing.T, abbrev string) *schema.Manager {
	t.Helper()
	mgr, err := schema.Open(filepath.Join(t.TempDir(), "bible.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateAppSchema(); err != nil {
		t.Fatal(err)
	}
	if err := mgr.DropTables(abbrev); err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateTables(abbrev); err != nil {
		t.Fatal(err)
	}

	return mgr
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
