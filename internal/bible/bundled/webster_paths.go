package bundled

import (
	"fmt"
	"os"
	"path/filepath"
)

const websterRel = "websters1828/dictionary.json"

var legacyDictionaryDBRelPaths = []string{
	"legacy/database/nativephp.sqlite",
	"legacy/storage/app/database/nativephp.sqlite",
}

// FindWebsterDictionaryJSON locates bundled dictionary.json (app resources or dev repo tree).
func FindWebsterDictionaryJSON(resourcesDir string) (string, error) {
	if resourcesDir != "" {
		p := filepath.Join(ResourcesRoot(resourcesDir), websterRel)
		if isFile(p) {
			return p, nil
		}
	}
	if p, ok := findUpwards(filepath.Join("extras", websterRel)); ok {
		return p, nil
	}
	return "", fmt.Errorf("webster dictionary JSON not found")
}

// FindLegacyDictionaryDB locates the reference NativePHP bible DB with imported Webster data.
func FindLegacyDictionaryDB() (string, error) {
	for _, rel := range legacyDictionaryDBRelPaths {
		if p, ok := findUpwards(rel); ok {
			return p, nil
		}
	}
	return "", fmt.Errorf("legacy dictionary database not found")
}

func findUpwards(rel string) (string, bool) {
	wd, err := os.Getwd()
	if err != nil {
		return "", false
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		candidate := filepath.Join(dir, rel)
		if isFile(candidate) {
			return candidate, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
	}
}

func isFile(path string) bool {
	st, err := os.Stat(path)
	return err == nil && !st.IsDir()
}
