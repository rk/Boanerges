package studyprint_test

import (
	"path/filepath"
	"strings"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/settings"
	"boanerges/internal/studyprint"
)

func studyPrintFixture(t *testing.T) string {
	t.Helper()
	root := filepath.Join("..", "..")
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	if err := bundled.ImportASV(
		db,
		filepath.Join(root, "extras", "bundled", "asv.sqlite"),
		filepath.Join(root, "extras", "sword", "mods.d", "asv.conf"),
	); err != nil {
		t.Skip(err)
	}
	return db
}

func TestBuildStudyHTMLGenesis(t *testing.T) {
	db := studyPrintFixture(t)

	html, err := studyprint.BuildStudyHTML(db, settings.StudySettings{
		TranslationID: "asv",
		BookID:        "gen",
		Chapter:       1,
	}, settings.ReadabilitySettings{FontSize: 16, LineHeight: 1.5, FontFamily: "serif"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "Genesis") || !strings.Contains(strings.ToLower(html), "beginning") {
		t.Fatalf("html = %q", html[:min(200, len(html))])
	}
	if !strings.Contains(html, `class="columns"`) {
		t.Fatalf("expected multi-column layout, html head = %q", html[:min(400, len(html))])
	}
}

func TestBuildStudyHTMLNotesColumnBlank(t *testing.T) {
	db := studyPrintFixture(t)

	html, err := studyprint.BuildStudyHTML(db, settings.StudySettings{
		TranslationID: "asv",
		BookID:        "gen",
		Chapter:       1,
		ColumnCount:   2,
		Columns:       []string{settings.ColumnNotes},
	}, settings.ReadabilitySettings{FontSize: 16, LineHeight: 1.5, FontFamily: "serif"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "column-lined-notes") {
		t.Fatalf("expected notes column, html snippet = %q", html[:min(500, len(html))])
	}
	if !strings.Contains(html, "lined-block") {
		t.Fatal("expected lined notes area when includeUserWork is false")
	}
}
