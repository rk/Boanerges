package progress

import "github.com/egoist/mygo"

// Step matches legacy App\Enums\TranslationInstallStep values.
type Step string

const (
	StepPending        Step = "pending"
	StepQueued         Step = "queued"
	StepCreatingSchema Step = "creating_schema"
	StepImporting      Step = "importing"
	StepIndexing       Step = "indexing"
	StepReady          Step = "ready"
	StepFailed         Step = "failed"
)

func (s Step) Percent() int {
	switch s {
	case StepReady:
		return 100
	case StepIndexing:
		return 85
	case StepImporting:
		return 50
	case StepCreatingSchema:
		return 30
	default:
		return 0
	}
}

// Payload matches legacy TranslationInstallProgress broadcast fields.
type Payload struct {
	Abbrev        string `json:"abbrev"`
	Step          string `json:"step"`
	Percent       int    `json:"percent"`
	Error         string `json:"error,omitempty"`
	InstallStatus string `json:"install_status,omitempty"`
}

// TranslationInstallProgress is emitted during translation installs.
var TranslationInstallProgress = mygo.NewEvent[Payload]("translationInstallProgress")

// Emitter broadcasts install progress to the UI.
type Emitter struct {
	win *mygo.Window
}

func NewEmitter() *Emitter {
	return &Emitter{}
}

func (e *Emitter) SetWindow(win *mygo.Window) {
	e.win = win
}

func (e *Emitter) Emit(payload Payload) {
	if e.win != nil {
		_ = TranslationInstallProgress.Emit(e.win, payload)
	}
}
