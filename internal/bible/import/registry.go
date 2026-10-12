package bibleimport

import (
	"fmt"

	"boanerges/internal/bible/catalog"
	"boanerges/internal/bible/import/swordfmt"
	"boanerges/internal/bible/import/usfm"
)

type FormatImporter interface {
	Format() catalog.ImportFormat
}

type Registry struct {
	sword *swordfmt.Importer
	usfm  *usfm.Importer
}

func NewRegistry(sword *swordfmt.Importer, usfm *usfm.Importer) *Registry {
	return &Registry{sword: sword, usfm: usfm}
}

func (r *Registry) Sword() *swordfmt.Importer { return r.sword }

func (r *Registry) Usfm() *usfm.Importer { return r.usfm }

func (r *Registry) Format(format catalog.ImportFormat) (FormatImporter, error) {
	switch format {
	case catalog.FormatSword:
		return r.sword, nil
	case catalog.FormatUsfm:
		return r.usfm, nil
	default:
		return nil, fmt.Errorf("unsupported import format %q", format)
	}
}
