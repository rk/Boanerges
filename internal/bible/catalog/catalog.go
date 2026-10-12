package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// Catalog loads translation catalog entries from JSON.
type Catalog struct {
	path    string
	entries []Entry
}

func Load(path string) (*Catalog, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rows []map[string]any
	if err := json.Unmarshal(raw, &rows); err != nil {
		return nil, err
	}

	entries := make([]Entry, 0, len(rows))
	for _, row := range rows {
		entry, err := entryFromRaw(row)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return &Catalog{path: path, entries: entries}, nil
}

func (c *Catalog) All() []Entry {
	out := make([]Entry, len(c.entries))
	copy(out, c.entries)

	return out
}

func (c *Catalog) Find(moduleKey string) (Entry, error) {
	for _, e := range c.entries {
		if strings.EqualFold(e.Short, moduleKey) {
			return e, nil
		}
	}

	return Entry{}, fmt.Errorf("translation %q is not in the catalog", moduleKey)
}

func (c *Catalog) FindByID(id string) (Entry, error) {
	id = strings.ToLower(id)
	for _, e := range c.entries {
		if e.ID() == id {
			return e, nil
		}
	}

	return Entry{}, fmt.Errorf("translation %q is not in the catalog", id)
}
