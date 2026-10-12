package sword

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseASVConf(t *testing.T) {
	path := filepath.Join("..", "..", "..", "extras", "sword", "mods.d", "asv.conf")
	meta, err := ParseConfFile(path)
	if err != nil {
		t.Skip("bundled ASV conf not present:", err)
	}

	if meta.Versification != "KJV" || meta.Format != "osis" || meta.VersionString != "2.0" {
		t.Fatalf("metadata = %+v", meta)
	}
	if meta.VersionDate != "2021-02-18" || meta.Copyright != "Public Domain" {
		t.Fatalf("dates/copyright = %+v", meta)
	}
	if meta.Source != "http://www.ebible.org/bible/asv/" {
		t.Fatalf("source = %q", meta.Source)
	}
	if !strings.Contains(meta.About, "American Standard Version") || strings.Contains(meta.About, `\par`) {
		t.Fatalf("about = %q", meta.About)
	}
}
