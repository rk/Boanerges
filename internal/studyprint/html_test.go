package studyprint_test

import (
	"path/filepath"
	"strings"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/settings"
	"boanerges/internal/studyprint"
)

func TestBuildStudyHTMLGenesis(t *testing.T) {
	root := filepath.Join("..", "..")
	db := filepath.Join(t.TempDir(), "bible.sqlite")
	if err := bundled.ImportASV(
		db,
		filepath.Join(root, "extras", "bundled", "asv.sqlite"),
		filepath.Join(root, "extras", "sword", "mods.d", "asv.conf"),
	); err != nil {
		t.Skip(err)
	}

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
}
