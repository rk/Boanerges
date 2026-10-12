package sword

import (
	"os"
	"path/filepath"
	"strings"
)

// Config holds parsed SWORD module .conf fields.
type Config struct {
	ModuleKey   string
	Abbrev      string
	Description string
	DataPath    string
	ModDrv      string
	SourceType  string
	Encoding    string
	Compress    string
	BlockType   string
	Versification string
	About       string
}

// Metadata is catalog-facing fields from a conf file.
type Metadata struct {
	Name             string
	Format           string
	Versification    string
	About            string
	VersionString    string
	VersionDate      string
	Copyright        string
	CopyrightContact string
	Source           string
}

// ParseConfFile reads and parses a SWORD .conf file.
func ParseConfFile(path string) (Config, Metadata, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Config{}, Metadata{}, err
	}

	cfg, meta := ParseConf(string(bytes))
	cfg.ModuleKey = strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	return cfg, meta, nil
}

// ParseConf parses SWORD .conf contents.
func ParseConf(raw string) (Config, Metadata) {
	data := parseRawConf(raw)

	about := data["About"]
	if about != "" {
		about = strings.TrimSpace(strings.ReplaceAll(about, `\par`, "\n\n"))
	}

	format := strings.ToLower(data["SourceType"])

	versionDate := data["SwordVersionDate"]
	if versionDate == "" {
		versionDate = data["RevDate"]
	}

	copyright := data["DistributionLicense"]
	if copyright == "" {
		copyright = data["Copyright"]
	}

	contact := data["CopyrightHolder"]
	if contact == "" {
		contact = data["CopyrightContact"]
	}

	cfg := Config{
		Abbrev:        data["Abbreviation"],
		Description:   data["Description"],
		DataPath:      data["DataPath"],
		ModDrv:        data["ModDrv"],
		SourceType:    format,
		Encoding:      data["Encoding"],
		Compress:      data["CompressType"],
		BlockType:     data["BlockType"],
		Versification: data["Versification"],
		About:         about,
	}

	meta := Metadata{
		Name:             data["Description"],
		Format:           format,
		Versification:    data["Versification"],
		About:            about,
		VersionString:    data["Version"],
		VersionDate:      versionDate,
		Copyright:        copyright,
		CopyrightContact: contact,
		Source:           data["TextSource"],
	}

	return cfg, meta
}

func parseRawConf(raw string) map[string]string {
	data := map[string]string{}
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, "[") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		data[strings.TrimSpace(key)] = strings.TrimSpace(value)
	}

	return data
}
