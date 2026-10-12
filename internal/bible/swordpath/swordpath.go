package swordpath

import (
	"os"
	"path/filepath"
	"strings"
)

// Roots returns bundled then local SWORD install roots.
func Roots(bundledRoot, localRoot string) []string {
	out := []string{}
	if bundledRoot != "" {
		out = append(out, bundledRoot)
	}
	if localRoot != "" {
		_ = os.MkdirAll(localRoot, 0o755)
		out = append(out, localRoot)
	}

	return out
}

// FindConf locates mods.d/{module}.conf under roots (case-insensitive module key).
func FindConf(roots []string, moduleKey string) (string, error) {
	for _, root := range roots {
		confDir := filepath.Join(root, "mods.d")
		entries, err := os.ReadDir(confDir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".conf") {
				continue
			}
			key := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
			if strings.EqualFold(key, moduleKey) {
				return filepath.Join(confDir, e.Name()), nil
			}
		}
	}

	return "", os.ErrNotExist
}

// ModuleInstalled reports whether moduleKey has a conf under roots.
func ModuleInstalled(roots []string, moduleKey string) bool {
	_, err := FindConf(roots, moduleKey)

	return err == nil
}
