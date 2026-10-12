package bundled

import (
	"path/filepath"
)

// Bootstrap copies bundled ASV into userDB if not already present.
func Bootstrap(userDB, resourcesDir string) error {
	return ImportASV(
		userDB,
		filepath.Join(resourcesDir, SQLiteFile),
		filepath.Join(resourcesDir, ConfFile),
	)
}
