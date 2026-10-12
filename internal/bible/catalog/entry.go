package catalog

import (
	"fmt"
	"strings"
)

// ImportFormat matches legacy CatalogImportFormat.
type ImportFormat string

const (
	FormatSword      ImportFormat = "sword"
	FormatUsfm       ImportFormat = "usfm"
	FormatAccordance ImportFormat = "accordance" // ponytail: not ported; catalog has zero active entries
)

// Entry is one row from extras/translations.json.
type Entry struct {
	Short        string       `json:"short"`
	Name         string       `json:"name"`
	URL          string       `json:"url"`
	About        string       `json:"about"`
	ImportAs     ImportFormat `json:"-"`
	MarkupFormat string       `json:"-"`
}

func (e Entry) ID() string {
	return strings.ToLower(e.Short)
}

func (e Entry) Module() string {
	return e.Short
}

func entryFromRaw(raw map[string]any) (Entry, error) {
	short, _ := raw["short"].(string)
	name, _ := raw["name"].(string)
	if short == "" || name == "" {
		return Entry{}, fmt.Errorf("catalog entry missing short or name")
	}

	importAs := FormatSword
	if v, ok := raw["import_as"].(string); ok && v != "" {
		importAs = ImportFormat(strings.ToLower(v))
		if importAs != FormatSword && importAs != FormatUsfm && importAs != FormatAccordance {
			importAs = FormatSword
		}
	} else if v, ok := raw["format"].(string); ok {
		switch strings.ToLower(v) {
		case "usfm":
			importAs = FormatUsfm
		case "sword":
			importAs = FormatSword
		}
	}

	url := ""
	if v, ok := raw["url"].(string); ok {
		url = v
	}
	if url == "" {
		url = fmt.Sprintf("https://crosswire.org/ftpmirror/pub/sword/packages/rawzip/%s", short)
	}

	about := ""
	if v, ok := raw["about"].(string); ok {
		about = v
	}
	if about == "" {
		about = fmt.Sprintf("https://crosswire.org/sword/modules/ModInfo.jsp?modName=%s", short)
	}

	markup := ""
	if v, ok := raw["markup_format"].(string); ok {
		markup = strings.ToLower(v)
	}

	return Entry{
		Short:        short,
		Name:         name,
		URL:          url,
		About:        about,
		ImportAs:     importAs,
		MarkupFormat: markup,
	}, nil
}
