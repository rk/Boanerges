package install

import (
	"strings"

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

	if existing != nil && existing.InstallStatus != translations.StatusFailed && existing.InstallStatus != translations.StatusPending {
		return nil, conflict("The " + entry.Short + " translation is already being installed.")
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
	} else {
		rec = existing
		_ = i.store.UpdateCatalogMeta(rec.ID, entry.Name, entry.MarkupFormat, entry.About)
	}

	if err != nil {
		return nil, err
	}

	_ = i.store.UpdateProgress(rec.ID, translations.StatusPending, "queued", nil)

	go func(id int64) {
		_ = i.pipeline.Run(id)
	}(rec.ID)

	return i.store.FindByID(rec.ID)
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
