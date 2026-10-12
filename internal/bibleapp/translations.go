package bibleapp

import (
	"path/filepath"

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

	"github.com/egoist/mygo"
)

// TranslationService exposes catalog install/uninstall to the frontend.
type TranslationService struct {
	catalog   *catalog.Catalog
	registry  *registry.Registry
	installer *install.Installer
	store     *translations.Store
}

func NewTranslationService(dbPath string, boot *bootstrap.Service) (*TranslationService, error) {
	mgr, err := schema.Open(dbPath)
	if err != nil {
		return nil, err
	}
	if err := mgr.CreateAppSchema(); err != nil {
		mgr.Close()
		return nil, err
	}

	resDir, err := mygo.App.Path(mygo.PathResources)
	if err != nil {
		mgr.Close()
		return nil, err
	}
	dataDir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		mgr.Close()
		return nil, err
	}

	cat, err := catalog.Load(filepath.Join(resDir, "translations.json"))
	if err != nil {
		mgr.Close()
		return nil, err
	}

	store := translations.NewStore(mgr)
	reg := registry.New(store)

	bundledSword := filepath.Join(resDir, "sword")
	localSword := filepath.Join(dataDir, "sword")
	roots := swordpath.Roots(bundledSword, localSword)

	swordImp := swordfmt.NewImporter(mgr, nil, roots)
	usfmImp := usfm.NewImporter(mgr, nil)
	importers := bibleimport.NewRegistry(swordImp, usfmImp)

	tmpDir := filepath.Join(dataDir, "tmp")
	pipeline := install.NewPipeline(cat, store, mgr, importers, bundledSword, localSword, tmpDir, boot.Emitter())
	installer := install.NewInstaller(cat, store, reg, mgr, pipeline)

	return &TranslationService{
		catalog:   cat,
		registry:  reg,
		installer: installer,
		store:     store,
	}, nil
}

type catalogResponse struct {
	Translations []translations.CatalogItem `json:"translations"`
}

type installedResponse struct {
	Translations []translations.Config `json:"translations"`
}

type installResponse struct {
	Translation translations.Config `json:"translation"`
}

// Catalog lists catalog entries with install status (legacy bible.translations.catalog).
func (s *TranslationService) Catalog() catalogResponse {
	items := make([]translations.CatalogItem, 0, len(s.catalog.All()))
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
			ID:            entry.ID(),
			Module:        entry.Module(),
			Name:          entry.Name,
			Abbrev:        entry.Short,
			Installed:     installed,
			Bundled:       s.registry.IsBundled(entry.Short),
			About:         entry.About,
			InstallStatus: status,
		})
	}

	return catalogResponse{Translations: items}
}

// Installed lists ready translations (legacy bible.translations.index).
func (s *TranslationService) Installed() installedResponse {
	recs, err := s.registry.Ready()
	if err != nil {
		return installedResponse{Translations: []translations.Config{}}
	}
	out := make([]translations.Config, 0, len(recs))
	for _, rec := range recs {
		out = append(out, s.registry.ToConfig(rec))
	}

	return installedResponse{Translations: out}
}

// Install queues a catalog translation install (legacy POST install).
func (s *TranslationService) Install(module string) (installResponse, error) {
	rec, err := s.installer.Install(module)
	if err != nil {
		return installResponse{}, err
	}

	return installResponse{Translation: s.registry.ToConfig(*rec)}, nil
}

// Uninstall removes an installed translation.
func (s *TranslationService) Uninstall(module string) error {
	return s.installer.Uninstall(module)
}

// InstallStatus returns in-progress install state for a module.
func (s *TranslationService) InstallStatus(module string) (map[string]any, error) {
	return s.installer.InstallStatus(module)
}
