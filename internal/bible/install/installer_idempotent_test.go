package install_test

import (
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/translations"
)

func TestInstallIdempotentWhilePending(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	userDB := filepath.Join(t.TempDir(), "bible.sqlite")
	if err := bundled.ImportASV(
		userDB,
		filepath.Join(root, "extras", "bundled", "asv.sqlite"),
		filepath.Join(root, "extras", "sword", "mods.d", "asv.conf"),
	); err != nil {
		t.Skip("bundled ASV:", err)
	}

	stack := newTestStack(t, userDB, root)
	first, err := stack.installer.Install("KJV")
	if err != nil {
		t.Fatal(err)
	}
	if first == nil || first.InstallStatus != translations.StatusPending {
		t.Fatalf("first = %+v", first)
	}

	second, err := stack.installer.Install("KJV")
	if err != nil {
		t.Fatal(err)
	}
	if second == nil || second.ID != first.ID {
		t.Fatalf("second = %+v", second)
	}
}
