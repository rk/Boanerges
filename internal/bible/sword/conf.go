package sword

import (
	"os"
	"strings"
)

// Metadata is parsed SWORD module .conf fields used by the app catalog.
type Metadata struct {
	Name              string `json:"name"`
	Format            string `json:"format"`
	Versification     string `json:"versification"`
	About             string `json:"about"`
	VersionString     string `json:"versionString"`
	VersionDate       string `json:"versionDate"`
	Copyright         string `json:"copyright"`
	CopyrightContact  string `json:"copyrightContact"`
	Source            string `json:"source"`
}

// ParseConfFile reads a SWORD .conf file into Metadata.
func ParseConfFile(path string) (Metadata, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return Metadata{}, err
	}

	return ParseConf(string(bytes)), nil
}

// ParseConf parses SWORD .conf contents.
func ParseConf(raw string) Metadata {
	data := parseRawConf(raw)

	about := data["About"]
	if about != "" {
		about = strings.TrimSpace(strings.ReplaceAll(about, `\par`, "\n\n"))
	}

	format := data["SourceType"]
	if format != "" {
		format = strings.ToLower(format)
	}

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

	return Metadata{
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
