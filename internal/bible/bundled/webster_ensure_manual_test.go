package bundled_test

import (
	"os"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/schema"
)

func TestEnsureWebster1828ManualEnv(t *testing.T) {
	res := os.Getenv("BOANERGES_RESOURCES")
	dbPath := os.Getenv("BOANERGES_DB")
	if res == "" || dbPath == "" {
		t.Skip("set BOANERGES_RESOURCES and BOANERGES_DB")
	}
	mgr, err := schema.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.CreateAppSchema(); err != nil {
		t.Fatal(err)
	}
	mgr.Close()

	if err := bundled.EnsureWebster1828(dbPath, res); err != nil {
		t.Fatal(err)
	}
	n := countDictionaryEntries(t, dbPath)
	if n == 0 {
		t.Fatal("no dictionary entries after ensure")
	}
	t.Logf("entries=%d json=%s", n, filepath.Join(bundled.ResourcesRoot(res), "websters1828/dictionary.json"))
}
