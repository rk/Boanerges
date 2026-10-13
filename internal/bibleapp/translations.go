package bibleapp

import (
	"path/filepath"
	"sync"

	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/bible/bundled"
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
	dbPath string
	boot   *bootstrap.Service

	once sync.Once
	err  error

	catalog   *catalog.Catalog
	registry  *registry.Registry
	installer *install.Installer
	store     *translations.Store
	mgr       *schema.Manager
}

func NewTranslationService(dbPath string, boot *bootstrap.Service) *TranslationService {
	return &TranslationService{dbPath: dbPath, boot: boot}
}

func (s *TranslationService) setup() error {
	s.once.Do(func() {
		s.mgr, s.err = schema.Open(s.dbPath)
		if s.err != nil {
			return
		}
		if s.err = s.mgr.CreateAppSchema(); s.err != nil {
			return
		}

		resDir, err := mygo.App.Path(mygo.PathResources)
		if err != nil {
			s.err = err
			return
		}
		dataDir, err := mygo.App.Path(mygo.PathUserData)
		if err != nil {
			s.err = err
			return
		}

		extras := bundled.ResourcesRoot(resDir)
		s.catalog, s.err = catalog.Load(filepath.Join(extras, "translations.json"))
		if s.err != nil {
			return
		}

		s.store = translations.NewStore(s.mgr)
		s.registry = registry.New(s.store)

		bundledSword := filepath.Join(extras, "sword")
		localSword := filepath.Join(dataDir, "sword")
		roots := swordpath.Roots(bundledSword, localSword)

		importers := bibleimport.NewRegistry(
			swordfmt.NewImporter(s.mgr, nil, roots),
			usfm.NewImporter(s.mgr, nil),
		)
		tmpDir := filepath.Join(dataDir, "tmp")
		pipeline := install.NewPipeline(s.catalog, s.store, s.mgr, importers, bundledSword, localSword, tmpDir, s.boot.Emitter())
		s.installer = install.NewInstaller(s.catalog, s.store, s.registry, s.mgr, pipeline)
	})

	return s.err
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
func (s *TranslationService) Catalog() (catalogResponse, error) {
	if err := s.setup(); err != nil {
		return catalogResponse{}, err
	}

	items := make([]translations.CatalogItem, 0, len(s.catalog.All()))
	for _, entry := range s.catalog.All() {
		rec, _ := s.store.FindByAbbrev(entry.ID())
		var status *string
		installed := false
		var step *string
		percent := 0
		if rec != nil {
			installed = rec.IsReady()
			if rec.InstallStatus != "" {
				st := rec.InstallStatus
				status = &st
			}
			if rec.InstallStep != "" {
				s := rec.InstallStep
				step = &s
			}
			if !rec.IsReady() {
				percent = translations.ProgressPercent(rec.InstallStatus, rec.InstallStep)
			}
		}
		items = append(items, translations.CatalogItem{
			ID: entry.ID(), Module: entry.Module(), Name: entry.Name, Abbrev: entry.Short,
			Installed: installed, Bundled: s.registry.IsBundled(entry.Short),
			About: entry.About, InstallStatus: status, InstallStep: step, InstallPercent: percent,
		})
	}

	return catalogResponse{Translations: items}, nil
}

// Installed lists ready translations (legacy bible.translations.index).
func (s *TranslationService) Installed() (installedResponse, error) {
	if err := s.setup(); err != nil {
		return installedResponse{}, err
	}
	recs, err := s.registry.Ready()
	if err != nil {
		return installedResponse{}, err
	}
	out := make([]translations.Config, 0, len(recs))
	for _, rec := range recs {
		out = append(out, s.registry.ToConfig(rec))
	}

	return installedResponse{Translations: out}, nil
}

// Install queues a catalog translation install (legacy POST install).
func (s *TranslationService) Install(module string) (installResponse, error) {
	if err := s.setup(); err != nil {
		return installResponse{}, err
	}
	rec, err := s.installer.Install(module)
	if err != nil {
		return installResponse{}, err
	}

	return installResponse{Translation: s.registry.ToConfig(*rec)}, nil
}

// Uninstall removes an installed translation.
func (s *TranslationService) Uninstall(module string) error {
	if err := s.setup(); err != nil {
		return err
	}

	return s.installer.Uninstall(module)
}

// InstallStatus returns in-progress install state for a module.
func (s *TranslationService) InstallStatus(module string) (map[string]any, error) {
	if err := s.setup(); err != nil {
		return nil, err
	}

	return s.installer.InstallStatus(module)
}

// ResumePendingInstalls re-queues stuck or interrupted catalog installs after app launch.
func (s *TranslationService) ResumePendingInstalls() error {
	if err := s.setup(); err != nil {
		return err
	}
	s.installer.ResumePending()

	return nil
}
