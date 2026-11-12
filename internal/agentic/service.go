//go:build agentic

package agentic

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"boanerges/internal/bibleapp"

	"github.com/egoist/mygo"
)

// Agentic exposes Studio APIs to the webview (agentic builds only).
type Agentic struct {
	bible  *bibleapp.BibleService
	dict   *bibleapp.DictionaryService
	runner *runnerRegistry
}

func NewAgentic(bible *bibleapp.BibleService, dict *bibleapp.DictionaryService) *Agentic {
	return &Agentic{
		bible:  bible,
		dict:   dict,
		runner: newRunnerRegistry(bible, dict),
	}
}

// PingOllama checks the local Ollama server.
func (a *Agentic) PingOllama() (OllamaHealth, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return PingOllama(ctx)
}

// PickProjectFolder opens a native directory picker.
func (a *Agentic) PickProjectFolder() (string, error) {
	paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
		Title:     "Choose project folder",
		Directory: true,
	})
	if err != nil {
		return "", err
	}
	if len(paths) == 0 {
		return "", nil
	}
	return paths[0], nil
}

// InitProject creates layout and manifest for a new or existing folder.
func (a *Agentic) InitProject(folder string, kind string) error {
	if folder == "" {
		return fmt.Errorf("folder required")
	}
	k := ProjectKind(kind)
	if k != ProjectDictionary && k != ProjectTopical {
		return fmt.Errorf("invalid kind")
	}
	return ensureProjectLayout(folder, k)
}

// LoadProjectSummary returns manifest, run state, and counts.
func (a *Agentic) LoadProjectSummary(folder string) (ProjectSummary, error) {
	if folder == "" {
		return ProjectSummary{}, fmt.Errorf("folder required")
	}
	return loadSummary(folder)
}

// ListSources lists files in sources/.
func (a *Agentic) ListSources(folder string) ([]string, error) {
	return listSources(folder)
}

// ListEntries lists dictionary entry summaries.
func (a *Agentic) ListEntries(folder string) ([]EntrySummary, error) {
	return listEntrySummaries(folder)
}

// ListTopics lists topical summaries.
func (a *Agentic) ListTopics(folder string) ([]TopicSummary, error) {
	return listTopicSummaries(folder)
}

// GetRunState returns the current run checkpoint.
func (a *Agentic) GetRunState(folder string) (RunState, error) {
	return loadRun(folder)
}

// SetProjectModel updates the Ollama model name in project.json.
func (a *Agentic) SetProjectModel(folder, model string) error {
	m, err := loadManifest(folder)
	if err != nil {
		return err
	}
	m.Model = model
	return writeJSON(filepath.Join(folder, manifestName), m)
}

func (a *Agentic) requireOllama() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	h, err := PingOllama(ctx)
	if err != nil {
		return err
	}
	if !h.OK {
		if h.Error != "" {
			return fmt.Errorf("ollama: %s", h.Error)
		}
		return fmt.Errorf("ollama unreachable")
	}
	return nil
}

// AnalyzeSourceFormats refreshes per-file entry-detection rules (10k sample, parallel subagents).
func (a *Agentic) AnalyzeSourceFormats(folder string) error {
	if err := a.requireOllama(); err != nil {
		return err
	}
	m, err := loadManifest(folder)
	if err != nil {
		return err
	}
	if m.Kind != ProjectDictionary {
		return fmt.Errorf("not a dictionary project")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()
	return ensureSourceGuides(ctx, folder, m.Model)
}

// StartDictionaryRun begins distilling entries from sources/.
func (a *Agentic) StartDictionaryRun(folder, instruction string) error {
	if err := a.requireOllama(); err != nil {
		return err
	}
	m, err := loadManifest(folder)
	if err != nil {
		return err
	}
	if m.Kind != ProjectDictionary {
		return fmt.Errorf("not a dictionary project")
	}
	return a.runner.start(folder, "distill", instruction, m.Model, nil)
}

// StartTopicalFromReference seeds topics from a verse.
func (a *Agentic) StartTopicalFromReference(folder, bookID string, chapter, verse int, instruction string) error {
	if err := a.requireOllama(); err != nil {
		return err
	}
	m, err := loadManifest(folder)
	if err != nil {
		return err
	}
	if m.Kind != ProjectTopical {
		return fmt.Errorf("not a topical project")
	}
	params := map[string]any{"bookId": bookID, "chapter": chapter, "verse": verse}
	return a.runner.start(folder, "from_reference", instruction, m.Model, params)
}

// StartTopicalFromTopic researches a named topic across Scripture.
func (a *Agentic) StartTopicalFromTopic(folder, topic, instruction string) error {
	if err := a.requireOllama(); err != nil {
		return err
	}
	m, err := loadManifest(folder)
	if err != nil {
		return err
	}
	if m.Kind != ProjectTopical {
		return fmt.Errorf("not a topical project")
	}
	params := map[string]any{"topic": topic}
	return a.runner.start(folder, "from_topic", instruction, m.Model, params)
}

// PauseRun pauses the active orchestrator.
func (a *Agentic) PauseRun(folder string) error {
	return a.runner.pause(folder)
}

// ResumeRun resumes a paused orchestrator.
func (a *Agentic) ResumeRun(folder string) error {
	return a.runner.resume(folder)
}

// ReadRawRunLog returns the tail of run.raw.log in the project folder (debug).
func (a *Agentic) ReadRawRunLog(folder string, maxBytes int) (string, error) {
	if folder == "" {
		return "", fmt.Errorf("folder required")
	}
	return readRawRunLogTail(folder, maxBytes)
}
