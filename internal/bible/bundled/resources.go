package bundled

import "path/filepath"

// ResourcesRoot returns the bundled extras tree. mygo packs `resources: ["extras"]`
// under PathResources as an `extras/` subdirectory.
func ResourcesRoot(resourcesDir string) string {
	return filepath.Join(resourcesDir, "extras")
}
