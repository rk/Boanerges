package translations

// Record mirrors the translations table row.
type Record struct {
	ID            int64
	Abbrev        string
	Name          string
	Format        string
	About         string
	InstallStatus string
	InstallStep   string
	InstallError  string
	Bundled       bool
}

func (r Record) IsReady() bool {
	return r.InstallStatus == StatusReady
}

// Config is the legacy TranslationResource / TranslationConfig JSON shape.
type Config struct {
	ID            string `json:"id"`
	Module        string `json:"module"`
	Name          string `json:"name"`
	Abbrev        string `json:"abbrev"`
	Bundled       bool   `json:"bundled"`
	About         string `json:"about,omitempty"`
	InstallStatus string `json:"install_status,omitempty"`
	InstallStep   string `json:"install_step,omitempty"`
	InstallError  string `json:"install_error,omitempty"`
}

// CatalogItem is legacy CatalogTranslationResource shape.
type CatalogItem struct {
	ID            string  `json:"id"`
	Module        string  `json:"module"`
	Name          string  `json:"name"`
	Abbrev        string  `json:"abbrev"`
	Installed     bool    `json:"installed"`
	Bundled       bool    `json:"bundled"`
	About         string  `json:"about,omitempty"`
	InstallStatus *string `json:"install_status"`
}

const (
	StatusPending        = "pending"
	StatusDownloading    = "downloading"
	StatusCreatingSchema = "creating_schema"
	StatusImporting      = "importing"
	StatusVerifying      = "verifying"
	StatusIndexing       = "indexing"
	StatusReady          = "ready"
	StatusFailed         = "failed"
)
