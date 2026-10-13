package schema_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
)

func TestRebuildFtsBSBScale(t *testing.T) {
	if testing.Short() {
		t.Skip("bench skipped in -short")
	}
	if os.Getenv("BOANERGES_FTS_BENCH") == "" {
		t.Skip("set BOANERGES_FTS_BENCH=1 to run local fts bench")
	}
	src := filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Boanerges", "bible.sqlite")
	if _, err := os.Stat(src); err != nil {
		t.Skip("user bible.sqlite not present")
	}
	dst := filepath.Join(t.TempDir(), "bench.sqlite")
 copyFile(t, src, dst)

	mgr, err := schema.Open(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	start := time.Now()
	if err := mgr.RebuildFtsIndex("bsb", markup.DefaultFormatter()); err != nil {
		t.Fatal(err)
	}
	t.Logf("bsb fts rebuild: %s", time.Since(start).Round(time.Millisecond))
}

func copyFile(t *testing.T, src, dst string) {
	t.Helper()
	in, err := os.Open(src)
	if err != nil {
		t.Skip(err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		t.Fatal(err)
	}
}
