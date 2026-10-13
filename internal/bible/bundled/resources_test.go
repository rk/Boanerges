package bundled_test

import (
	"os"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
)

func TestResourcesRootPackagedLayout(t *testing.T) {
	root := t.TempDir()
	extras := filepath.Join(root, "extras", "bundled")
	if err := os.MkdirAll(extras, 0o755); err != nil {
		t.Fatal(err)
	}
	sqlite := filepath.Join(extras, "asv.sqlite")
	if err := os.WriteFile(sqlite, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(bundled.ResourcesRoot(root), bundled.SQLiteFile)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected bundled sqlite at %s: %v", path, err)
	}
}
