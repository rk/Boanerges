package install_test

import (
	"errors"
	"path/filepath"
	"testing"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/bible/catalog"
	bibleimport "boanerges/internal/bible/import"
	"boanerges/internal/bible/import/swordfmt"
	"boanerges/internal/bible/import/usfm"
	"boanerges/internal/bible/install"
	"boanerges/internal/bible/registry"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/swordpath"
	"boanerges/internal/bible/translations"
)

func TestBundledASVInstallGuards(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	userDB := filepath.Join(t.TempDir(), "bible.sqlite")
	bundledSQLite := filepath.Join(root, "extras", "bundled", "asv.sqlite")
	confPath := filepath.Join(root, "extras", "sword", "mods.d", "asv.conf")
	if err := bundled.ImportASV(userDB, bundledSQLite, confPath); err != nil {
		t.Skip("bundled ASV:", err)
	}

	stack := newTestStack(t, userDB, root)
	asv, err := stack.catalog.Find("ASV")
	if err != nil {
		t.Fatal(err)
	}
	items := stack.svcCatalog()
	var asvItem *translations.CatalogItem
	for i := range items {
		if items[i].Module == "ASV" {
			asvItem = &items[i]
			break
		}
	}
	if asvItem == nil || !asvItem.Installed || !asvItem.Bundled {
		t.Fatal("ASV should be installed bundled in catalog")
	}

	_, err = stack.installer.Install(asv.Short)
	if !errors.Is(err, install.ErrUnprocessable) {
		t.Fatalf("install bundled ASV: %v", err)
	}

	err = stack.installer.Uninstall(asv.Short)
	if !errors.Is(err, install.ErrUnprocessable) {
		t.Fatalf("uninstall bundled ASV: %v", err)
	}
}

func TestInstallConflictWhenReady(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	userDB := filepath.Join(t.TempDir(), "bible.sqlite")
	bundledSQLite := filepath.Join(root, "extras", "bundled", "asv.sqlite")
	confPath := filepath.Join(root, "extras", "sword", "mods.d", "asv.conf")
	if err := bundled.ImportASV(userDB, bundledSQLite, confPath); err != nil {
		t.Skip("bundled ASV:", err)
	}

	stack := newTestStack(t, userDB, root)
	_, err := stack.installer.Install("ASV")
	if !errors.Is(err, install.ErrUnprocessable) {
		t.Fatalf("expected unprocessable, got %v", err)
	}
}

type testStack struct {
	catalog   *catalog.Catalog
	installer *install.Installer
	store     *translations.Store
}

func (s *testStack) svcCatalog() []translations.CatalogItem {
	items := make([]translations.CatalogItem, 0)
	for _, entry := range s.catalog.All() {
		rec, _ := s.store.FindByAbbrev(entry.ID())
		var status *string
		installed := false
		if rec != nil {
			installed = rec.IsReady()
			if rec.InstallStatus != "" {
				st := rec.InstallStatus
				status = &st
			}
		}
		items = append(items, translations.CatalogItem{
			ID: entry.ID(), Module: entry.Module(), Name: entry.Name, Abbrev: entry.Short,
			Installed: installed, Bundled: entry.Short == "ASV" || entry.ID() == "asv",
			About: entry.About, InstallStatus: status,
		})
	}

	return items
}

func newTestStack(t *testing.T, userDB, root string) *testStack {
	t.Helper()
	mgr, err := schema.Open(userDB)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { mgr.Close() })

	cat, err := catalog.Load(filepath.Join(root, "extras", "translations.json"))
	if err != nil {
		t.Fatal(err)
	}

	store := translations.NewStore(mgr)
	reg := registry.New(store)
	bundledSword := filepath.Join(root, "extras", "sword")
	localSword := filepath.Join(t.TempDir(), "sword")
	roots := swordpath.Roots(bundledSword, localSword)
	importers := bibleimport.NewRegistry(
		swordfmt.NewImporter(mgr, nil, roots),
		usfm.NewImporter(mgr, nil),
	)
	pipeline := install.NewPipeline(cat, store, mgr, importers, bundledSword, localSword, t.TempDir(), bootstrap.NewService().Emitter())
	installer := install.NewInstaller(cat, store, reg, mgr, pipeline)

	return &testStack{catalog: cat, installer: installer, store: store}
}
