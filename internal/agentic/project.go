//go:build agentic

package agentic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	manifestName = "project.json"
	runName      = "run.json"
)

func projectDirs(root string) (sources, entries, topics string) {
	return filepath.Join(root, "sources"), filepath.Join(root, "entries"), filepath.Join(root, "topics")
}

func ensureProjectLayout(root string, kind ProjectKind) error {
	sources, entries, topics := projectDirs(root)
	for _, dir := range []string{root, sources, entries, topics} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	manifestPath := filepath.Join(root, manifestName)
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		m := ProjectManifest{Kind: kind, Model: "llama3.2"}
		return writeJSON(manifestPath, m)
	}
	return nil
}

func loadManifest(root string) (ProjectManifest, error) {
	var m ProjectManifest
	if err := readJSON(filepath.Join(root, manifestName), &m); err != nil {
		return m, err
	}
	if m.Model == "" {
		m.Model = "llama3.2"
	}
	return m, nil
}

func defaultRun() RunState {
	return RunState{Status: RunIdle, Steps: []RunStep{}, UpdatedAt: time.Now().UTC().Format(time.RFC3339)}
}

func loadRun(root string) (RunState, error) {
	var r RunState
	path := filepath.Join(root, runName)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return defaultRun(), nil
	}
	if err := readJSON(path, &r); err != nil {
		return r, err
	}
	if r.Steps == nil {
		r.Steps = []RunStep{}
	}
	return r, nil
}

func saveRun(root string, r RunState) error {
	r.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	return writeJSON(filepath.Join(root, runName), r)
}

func appendStep(root, role, summary string) error {
	r, err := loadRun(root)
	if err != nil {
		return err
	}
	r.Steps = append(r.Steps, RunStep{
		At:      time.Now().UTC().Format(time.RFC3339),
		Role:    role,
		Summary: summary,
	})
	// ponytail: cap log at 500 steps; upgrade path is rotate to run.old.json
	if len(r.Steps) > 500 {
		r.Steps = r.Steps[len(r.Steps)-500:]
	}
	return saveRun(root, r)
}

func countFiles(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	n := 0
	for _, e := range entries {
		if !e.IsDir() {
			n++
		}
	}
	return n
}

func loadSummary(root string) (ProjectSummary, error) {
	manifest, err := loadManifest(root)
	if err != nil {
		return ProjectSummary{}, err
	}
	run, err := loadRun(root)
	if err != nil {
		return ProjectSummary{}, err
	}
	sources, entriesDir, topicsDir := projectDirs(root)
	return ProjectSummary{
		Folder:      root,
		Manifest:    manifest,
		Run:         run,
		SourceCount: countFiles(sources),
		EntryCount:  countFiles(entriesDir),
		TopicCount:  countFiles(topicsDir),
	}, nil
}

func listSources(root string) ([]string, error) {
	sources, _, _ := projectDirs(root)
	entries, err := os.ReadDir(sources)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() {
			out = append(out, e.Name())
		}
	}
	return out, nil
}

func readSourceChunk(root, name string, offset, limit int) (string, error) {
	if strings.Contains(name, "..") || strings.ContainsAny(name, `/\`) {
		return "", fmt.Errorf("invalid source name")
	}
	sources, _, _ := projectDirs(root)
	data, err := os.ReadFile(filepath.Join(sources, name))
	if err != nil {
		return "", err
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= len(data) {
		return "", nil
	}
	end := offset + limit
	if limit <= 0 || end > len(data) {
		end = len(data)
	}
	return string(data[offset:end]), nil
}

func listEntrySummaries(root string) ([]EntrySummary, error) {
	_, entriesDir, _ := projectDirs(root)
	entries, err := os.ReadDir(entriesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []EntrySummary{}, nil
		}
		return nil, err
	}
	out := make([]EntrySummary, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var doc DictEntry
		if err := readJSON(filepath.Join(entriesDir, e.Name()), &doc); err != nil {
			continue
		}
		hw := doc.Headword
		if hw == "" {
			hw = strings.TrimSuffix(e.Name(), ".json")
		}
		out = append(out, EntrySummary{Headword: hw, File: e.Name()})
	}
	return out, nil
}

func listTopicSummaries(root string) ([]TopicSummary, error) {
	_, _, topicsDir := projectDirs(root)
	entries, err := os.ReadDir(topicsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []TopicSummary{}, nil
		}
		return nil, err
	}
	out := make([]TopicSummary, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		var doc TopicDoc
		if err := readJSON(filepath.Join(topicsDir, e.Name()), &doc); err != nil {
			continue
		}
		title := doc.Title
		if title == "" {
			title = doc.Slug
		}
		slug := doc.Slug
		if slug == "" {
			slug = strings.TrimSuffix(e.Name(), ".json")
		}
		out = append(out, TopicSummary{Slug: slug, Title: title})
	}
	return out, nil
}

func readEntry(root, headword string) (DictEntry, error) {
	hw := strings.TrimSpace(headword)
	if hw == "" {
		return DictEntry{}, fmt.Errorf("headword required")
	}
	_, entriesDir, _ := projectDirs(root)
	path := filepath.Join(entriesDir, slugify(hw)+".json")
	var doc DictEntry
	if err := readJSON(path, &doc); err != nil {
		return DictEntry{}, err
	}
	return doc, nil
}

func upsertEntry(root string, entry DictEntry) error {
	hw := strings.TrimSpace(entry.Headword)
	if hw == "" {
		return fmt.Errorf("headword required")
	}
	slug := slugify(hw)
	_, entriesDir, _ := projectDirs(root)
	return writeJSON(filepath.Join(entriesDir, slug+".json"), entry)
}

func upsertTopic(root string, doc TopicDoc) error {
	slug := strings.TrimSpace(doc.Slug)
	if slug == "" {
		slug = slugify(doc.Title)
	}
	if slug == "" {
		return fmt.Errorf("topic slug or title required")
	}
	doc.Slug = slug
	_, _, topicsDir := projectDirs(root)
	return writeJSON(filepath.Join(topicsDir, slug+".json"), doc)
}

func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && b.Len() > 0 {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
