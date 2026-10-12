package registry

import (
	"strings"

	"boanerges/internal/bible/translations"
)

var bundledModules = []string{"ASV"}

// Registry tracks installed translations from the database.
type Registry struct {
	store *translations.Store
}

func New(store *translations.Store) *Registry {
	return &Registry{store: store}
}

func (r *Registry) IsBundled(moduleKey string) bool {
	for _, m := range bundledModules {
		if strings.EqualFold(m, moduleKey) {
			return true
		}
	}

	return false
}

func (r *Registry) IsInstalled(moduleKey string) (bool, error) {
	rec, err := r.store.FindByAbbrev(moduleKey)
	if err != nil {
		return false, err
	}

	return rec != nil && rec.IsReady(), nil
}

func (r *Registry) FindModel(abbrev string) (*translations.Record, error) {
	return r.store.FindByAbbrev(abbrev)
}

func (r *Registry) Ready() ([]translations.Record, error) {
	return r.store.ListReady()
}

func (r *Registry) ToConfig(rec translations.Record) translations.Config {
	module := strings.ToUpper(rec.Abbrev)

	return translations.Config{
		ID:            rec.Abbrev,
		Module:        module,
		Name:          rec.Name,
		Abbrev:        strings.ToUpper(rec.Abbrev),
		Bundled:       rec.Bundled,
		About:         rec.About,
		InstallStatus: rec.InstallStatus,
		InstallStep:   rec.InstallStep,
		InstallError:  rec.InstallError,
	}
}
