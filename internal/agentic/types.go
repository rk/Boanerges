//go:build agentic

package agentic

type OllamaHealth struct {
	OK     bool     `json:"ok"`
	Error  string   `json:"error,omitempty"`
	Models []string `json:"models,omitempty"`
}

type ProjectKind string

const (
	ProjectDictionary ProjectKind = "dictionary"
	ProjectTopical    ProjectKind = "topical"
)

type ProjectManifest struct {
	Kind  ProjectKind `json:"kind"`
	Model string      `json:"model"`
}

type RunStatus string

const (
	RunIdle    RunStatus = "idle"
	RunRunning RunStatus = "running"
	RunPaused  RunStatus = "paused"
	RunFailed  RunStatus = "failed"
)

type RunStep struct {
	At      string `json:"at"`
	Role    string `json:"role"`
	Summary string `json:"summary"`
}

type RunState struct {
	Status      RunStatus          `json:"status"`
	Job         string             `json:"job"`
	Instruction string             `json:"instruction"`
	Params      any                `json:"params,omitempty"`
	Steps       []RunStep          `json:"steps"`
	Distill     *DistillCheckpoint `json:"distill,omitempty"`
	Error       string             `json:"error,omitempty"`
	UpdatedAt   string             `json:"updatedAt"`
}

type ProjectSummary struct {
	Folder       string        `json:"folder"`
	Manifest     ProjectManifest `json:"manifest"`
	Run          RunState      `json:"run"`
	SourceCount  int           `json:"sourceCount"`
	EntryCount   int           `json:"entryCount"`
	TopicCount   int           `json:"topicCount"`
}

type EntrySummary struct {
	Headword string `json:"headword"`
	File     string `json:"file"`
}

type TopicSummary struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

type DictEntry struct {
	Headword string       `json:"headword"`
	Modern   string       `json:"modern"`
	Senses   []DictSense  `json:"senses"`
}

type DictSense struct {
	Gloss      string `json:"gloss"`
	SourceFile string `json:"sourceFile"`
	SourceSpan string `json:"sourceSpan"`
}

type TopicDoc struct {
	Slug   string       `json:"slug"`
	Title  string       `json:"title"`
	Groups []TopicGroup `json:"groups"`
}

type TopicGroup struct {
	Label string       `json:"label"`
	Refs  []ScriptureRef `json:"refs"`
}

type ScriptureRef struct {
	BookID  string `json:"bookId"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
}
