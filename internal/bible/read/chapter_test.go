package read_test

import (
	"path/filepath"
	"strings"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/read"
	"boanerges/internal/bible/schema"
)

func TestGenesisChapterOne(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	bundledPath := filepath.Join(root, "extras", "bundled", "asv.sqlite")
	confPath := filepath.Join(root, "extras", "sword", "mods.d", "asv.conf")

	userDB := filepath.Join(t.TempDir(), "app.db")
	if err := bundled.ImportASV(userDB, bundledPath, confPath); err != nil {
		t.Skip("bundled ASV:", err)
	}

	mgr, err := schema.Open(userDB)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	reader := read.NewReader(mgr.DB(), mgr, markup.DefaultFormatter())
	ch, err := reader.Chapter("asv", "gen", 1)
	if err != nil {
		t.Fatal(err)
	}
	if ch.Book != "Genesis" || ch.Chapter != 1 || len(ch.Verses) != 31 {
		t.Fatalf("chapter = %+v", ch)
	}
	if !strings.Contains(strings.ToLower(ch.Verses[0].Text), "beginning") {
		t.Fatalf("verse 1 text = %q", ch.Verses[0].Text)
	}
}
