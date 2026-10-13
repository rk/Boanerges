package bundled

import (
	"path/filepath"
)

// Bootstrap copies bundled ASV into userDB if not already present.
func Bootstrap(userDB, resourcesDir string) error {
	root := ResourcesRoot(resourcesDir)
	return ImportASV(
		userDB,
		filepath.Join(root, SQLiteFile),
		filepath.Join(root, ConfFile),
	)
}
