package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/egoist/mygo"
)

// Repository persists nested settings keys in settings.json under user data.
type Repository struct {
	// Dir overrides user data directory (for tests). Empty uses mygo.PathUserData.
	Dir string

	mu sync.Mutex
}

func (r *Repository) Get(key string, defaults map[string]any) map[string]any {
	stored, ok := r.readKey(key)
	if !ok {
		return copyMap(defaults)
	}

	return mergeMaps(defaults, stored)
}

func (r *Repository) Update(key string, patch map[string]any, defaults map[string]any) map[string]any {
	merged := mergeMaps(r.Get(key, defaults), patch)
	r.writeKey(key, merged)

	return merged
}

func (r *Repository) readKey(key string) (map[string]any, bool) {
	data, err := r.readFile()
	if err != nil || data == nil {
		return nil, false
	}

	raw, ok := data[key]
	if !ok {
		return nil, false
	}

	stored, ok := raw.(map[string]any)
	if !ok {
		return nil, false
	}

	return stored, true
}

func (r *Repository) writeKey(key string, value map[string]any) {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, _ := r.readFileUnlocked()
	if data == nil {
		data = map[string]any{}
	}
	data[key] = value
	_ = r.writeFileUnlocked(data)
}

func (r *Repository) readFile() (map[string]any, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return r.readFileUnlocked()
}

func (r *Repository) readFileUnlocked() (map[string]any, error) {
	path, err := r.filePath()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return map[string]any{}, nil
	}
	if err != nil {
		return nil, err
	}

	var data map[string]any
	if err := json.Unmarshal(bytes, &data); err != nil {
		return map[string]any{}, nil
	}
	if data == nil {
		return map[string]any{}, nil
	}

	return data, nil
}

func (r *Repository) writeFileUnlocked(data map[string]any) error {
	path, err := r.filePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, encoded, 0o644)
}

func (r *Repository) filePath() (string, error) {
	dir := r.Dir
	if dir == "" {
		var err error
		dir, err = mygo.App.Path(mygo.PathUserData)
		if err != nil {
			return "", err
		}
	}

	return filepath.Join(dir, "settings.json"), nil
}

func mergeMaps(base map[string]any, overlay map[string]any) map[string]any {
	out := copyMap(base)
	for k, v := range overlay {
		out[k] = v
	}

	return out
}

func copyMap(in map[string]any) map[string]any {
	out := make(map[string]any, len(in))
	for k, v := range in {
		out[k] = v
	}

	return out
}

func asString(m map[string]any, key string, fallback string) string {
	v, ok := m[key]
	if !ok {
		return fallback
	}
	s, ok := v.(string)
	if !ok {
		return fallback
	}

	return s
}

func asInt(m map[string]any, key string, fallback int) int {
	v, ok := m[key]
	if !ok {
		return fallback
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case int64:
		return int(n)
	default:
		return fallback
	}
}

func asBool(m map[string]any, key string, fallback bool) bool {
	v, ok := m[key]
	if !ok {
		return fallback
	}
	b, ok := v.(bool)
	if !ok {
		return fallback
	}

	return b
}

func asStringSlice(m map[string]any, key string) []string {
	v, ok := m[key]
	if !ok {
		return nil
	}
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if ok {
			out = append(out, s)
		}
	}

	return out
}

func validationError(msg string) error {
	return fmt.Errorf("%s", msg)
}
