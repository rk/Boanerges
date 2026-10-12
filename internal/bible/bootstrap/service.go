package bootstrap

import (
	"path/filepath"
	"sync"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/install"

	"github.com/egoist/mygo"
)

// ProgressPayload matches legacy TranslationInstallProgress broadcast fields.
type ProgressPayload struct {
	Abbrev        string `json:"abbrev"`
	Step          string `json:"step"`
	Percent       int    `json:"percent"`
	Error         string `json:"error,omitempty"`
	InstallStatus string `json:"install_status,omitempty"`
}

// TranslationInstallProgress is emitted during bundled ASV copy (welcome flow).
var TranslationInstallProgress = mygo.NewEvent[ProgressPayload]("translationInstallProgress")

// Service runs first-launch bundled imports.
type Service struct {
	mu     sync.Mutex
	status ProgressPayload
	done   bool
	win    *mygo.Window
}

func NewService() *Service {
	return &Service{
		status: ProgressPayload{
			Abbrev:        "asv",
			Step:          string(install.StepPending),
			Percent:       0,
			InstallStatus: "pending",
		},
	}
}

// ASVStatus returns current bundled ASV install progress.
func (s *Service) ASVStatus() ProgressPayload {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.status
}

// ASVReady reports whether bundled ASV tables are present.
func (s *Service) ASVReady() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.done
}

func (s *Service) setStatus(step install.Step, installStatus string, errMsg string) {
	s.mu.Lock()
	s.status = ProgressPayload{
		Abbrev:        "asv",
		Step:          string(step),
		Percent:       step.Percent(),
		InstallStatus: installStatus,
		Error:         errMsg,
	}
	if step == install.StepReady {
		s.done = true
		s.status.Percent = 100
	}
	payload := s.status
	s.mu.Unlock()

	if s.win != nil {
		_ = TranslationInstallProgress.Emit(s.win, payload)
	}
}

// SetWindow sets the window used for install progress events.
func (s *Service) SetWindow(win *mygo.Window) {
	s.mu.Lock()
	s.win = win
	s.mu.Unlock()
}

// RunASV copies bundled sqlite into userDB (idempotent).
func (s *Service) RunASV(userDB, resourcesDir string) error {
	s.setStatus(install.StepQueued, "pending", "")

	bundledSQLite := filepath.Join(resourcesDir, bundled.SQLiteFile)
	confPath := filepath.Join(resourcesDir, bundled.ConfFile)

	s.setStatus(install.StepCreatingSchema, "creating_schema", "")
	if err := bundled.ImportASV(userDB, bundledSQLite, confPath); err != nil {
		s.setStatus(install.StepFailed, "failed", err.Error())
		return err
	}

	s.setStatus(install.StepIndexing, "indexing", "")
	s.setStatus(install.StepReady, "ready", "")

	return nil
}
