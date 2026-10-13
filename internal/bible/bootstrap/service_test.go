package bootstrap_test

import (
	"os"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/bible/bundled"
)

func TestRunASVPackagedResourcesLayout(t *testing.T) {
	repo := filepath.Join("..", "..", "..")
	packagedResources := filepath.Join(repo, "build", "darwin-arm64", "Boanerges.app", "Contents", "Resources")
	if _, err := os.Stat(packagedResources); err != nil {
		t.Skip("built app resources not present:", err)
	}

	db := filepath.Join(t.TempDir(), "bible.sqlite")
	svc := bootstrap.NewService()
	if err := svc.RunASV(db, packagedResources); err != nil {
		t.Fatal(err)
	}
	if !svc.ASVReady() {
		t.Fatal("expected ASV ready")
	}
	if _, err := os.Stat(db); err != nil {
		t.Fatal(err)
	}

	root := bundled.ResourcesRoot(packagedResources)
	if _, err := os.Stat(filepath.Join(root, bundled.SQLiteFile)); err != nil {
		t.Fatalf("bundled sqlite missing at extras path: %v", err)
	}
}
