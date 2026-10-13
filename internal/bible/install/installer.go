package install

import (
	"strings"
	"sync"

	"boanerges/internal/bible/catalog"
	"boanerges/internal/bible/registry"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/translations"
)

type Installer struct {
	catalog  *catalog.Catalog
	store    *translations.Store
	registry *registry.Registry
	schema   *schema.Manager
	pipeline *Pipeline

	mu      sync.Mutex
	running map[int64]struct{}
	runMu   sync.Mutex // one import at a time (single writer on bible.sqlite)
}

func NewInstaller(
	cat *catalog.Catalog,
	store *translations.Store,
	reg *registry.Registry,
	mgr *schema.Manager,
	pipeline *Pipeline,
) *Installer {
	return &Installer{
		catalog:  cat,
		store:    store,
		registry: reg,
		schema:   mgr,
		pipeline: pipeline,
		running:  make(map[int64]struct{}),
	}
}

func (i *Installer) Install(moduleKey string) (*translations.Record, error) {
	entry, err := i.catalog.Find(moduleKey)
	if err != nil {
		return nil, notFound(err.Error())
	}

	if i.registry.IsBundled(entry.Short) {
		installed, err := i.registry.IsInstalled(entry.Short)
		if err != nil {
			return nil, err
		}
		if installed {
			return nil, unprocessable("The " + entry.Short + " translation is bundled with the app and cannot be installed separately.")
		}
	}

	existing, err := i.store.FindByAbbrev(entry.ID())
	if err != nil {
		return nil, err
	}

	if existing != nil && existing.IsReady() {
		return nil, conflict("The " + entry.Short + " translation is already installed.")
	}

	var rec *translations.Record
	if existing == nil {
		format := entry.MarkupFormat
		if format == "" && entry.ImportAs == catalog.FormatUsfm {
			format = "usfm"
		}
		rec, err = i.store.Create(translations.CreateParams{
			Abbrev:  entry.ID(),
			Name:    entry.Name,
			Format:  format,
			About:   entry.About,
			Bundled: false,
		})
		if err != nil {
			return nil, err
		}
	} else {
		rec = existing
		_ = i.store.UpdateCatalogMeta(rec.ID, entry.Name, entry.MarkupFormat, entry.About)
		if rec.InstallStatus == translations.StatusFailed {
			_ = i.store.UpdateProgress(rec.ID, translations.StatusPending, "queued", nil)
		}
	}

	_ = i.store.UpdateProgress(rec.ID, translations.StatusPending, "queued", nil)
	i.enqueue(rec.ID)

	return i.store.FindByID(rec.ID)
}

// ResumePending re-queues incomplete catalog installs (e.g. stuck on indexing after restart).
func (i *Installer) ResumePending() {
	recs, err := i.store.ListIncomplete()
	if err != nil {
		return
	}
	for _, rec := range recs {
		if i.registry.IsBundled(rec.Abbrev) {
			continue
		}
		i.enqueue(rec.ID)
	}
}

func (i *Installer) enqueue(translationID int64) {
	i.mu.Lock()
	if _, active := i.running[translationID]; active {
		i.mu.Unlock()
		return
	}
	i.running[translationID] = struct{}{}
	i.mu.Unlock()

	go func(id int64) {
		defer func() {
			i.mu.Lock()
			delete(i.running, id)
			i.mu.Unlock()
		}()
		i.runMu.Lock()
		defer i.runMu.Unlock()
		_ = i.pipeline.Run(id)
	}(translationID)
}

func (i *Installer) Uninstall(moduleKey string) error {
	entry, err := i.catalog.Find(moduleKey)
	if err != nil {
		return notFound(err.Error())
	}

	if i.registry.IsBundled(entry.Short) {
		return unprocessable("The " + entry.Short + " translation is bundled with the app and cannot be removed.")
	}

	rec, err := i.store.FindByAbbrev(entry.ID())
	if err != nil {
		return err
	}
	if rec == nil || !rec.IsReady() {
		return notFound("The " + entry.Short + " translation is not installed.")
	}

	a, err := i.schema.ValidateAbbrev(rec.Abbrev)
	if err != nil {
		return err
	}
	if err := i.schema.DropTables(a); err != nil {
		return err
	}

	return i.store.Delete(rec.ID)
}

func (i *Installer) InstallStatus(moduleKey string) (map[string]any, error) {
	entry, err := i.catalog.Find(moduleKey)
	if err != nil {
		if _, err2 := i.catalog.FindByID(strings.ToLower(moduleKey)); err2 != nil {
			return nil, notFound("Translation not found.")
		}
		entry, _ = i.catalog.FindByID(strings.ToLower(moduleKey))
	}

	rec, err := i.store.FindByAbbrev(entry.ID())
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, notFound("Translation not found.")
	}

	step := rec.InstallStep
	if step == "" {
		step = "pending"
	}

	return map[string]any{
		"abbrev":         rec.Abbrev,
		"install_status": rec.InstallStatus,
		"step":           step,
		"percent":        translations.ProgressPercent(rec.InstallStatus, step),
		"install_error":  rec.InstallError,
	}, nil
}
