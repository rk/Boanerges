package sword

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestOpenASVGenesisVerse(t *testing.T) {
	conf := filepath.Join("..", "extras", "sword", "mods.d", "asv.conf")
	mod, err := OpenModule(conf)
	if err != nil {
		t.Skip("bundled ASV module:", err)
	}

	text, err := mod.GetVerse("gen", 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text, "beginning") {
		t.Fatalf("Gen 1:1 = %q", text)
	}
}

func TestOpenASVJohn316(t *testing.T) {
	conf := filepath.Join("..", "extras", "sword", "mods.d", "asv.conf")
	mod, err := OpenModule(conf)
	if err != nil {
		t.Skip("bundled ASV module:", err)
	}

	text, err := mod.GetVerse("jhn", 3, 16)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.ToLower(text), "loved") {
		t.Fatalf("John 3:16 = %q", text)
	}
}

func TestParseASVConf(t *testing.T) {
	confPath := filepath.Join("..", "extras", "sword", "mods.d", "asv.conf")
	_, meta, err := ParseConfFile(confPath)
	if err != nil {
		t.Skip(err)
	}
	if meta.Versification != "KJV" || meta.Format != "osis" {
		t.Fatalf("meta = %+v", meta)
	}
}
