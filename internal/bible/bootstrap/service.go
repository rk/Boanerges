package bootstrap

import (
	"path/filepath"
	"sync"

	"boanerges/internal/bible/bundled"
	"boanerges/internal/bible/progress"

	"github.com/egoist/mygo"
)

// ProgressPayload is the install progress JSON shape (alias for mygo bindings).
type ProgressPayload = progress.Payload

// TranslationInstallProgress re-exports the mygo event for bindings.
var TranslationInstallProgress = progress.TranslationInstallProgress

// Service runs first-launch bundled imports.
type Service struct {
	mu      sync.Mutex
	status  progress.Payload
	done    bool
	emitter *progress.Emitter
}

func NewService() *Service {
	return &Service{
		emitter: progress.NewEmitter(),
		status: progress.Payload{
			Abbrev:        "asv",
			Step:          string(progress.StepPending),
			Percent:       0,
			InstallStatus: "pending",
		},
	}
}

// Emitter returns shared install progress broadcaster.
func (s *Service) Emitter() *progress.Emitter {
	return s.emitter
}

// ASVStatus returns current bundled ASV install progress.
func (s *Service) ASVStatus() progress.Payload {
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

func (s *Service) setStatus(step progress.Step, installStatus string, errMsg string) {
	s.mu.Lock()
	s.status = progress.Payload{
		Abbrev:        "asv",
		Step:          string(step),
		Percent:       step.Percent(),
		InstallStatus: installStatus,
		Error:         errMsg,
	}
	if step == progress.StepReady {
		s.done = true
		s.status.Percent = 100
	}
	payload := s.status
	s.mu.Unlock()

	s.emitter.Emit(payload)
}

// SetWindow sets the window used for install progress events.
func (s *Service) SetWindow(win *mygo.Window) {
	s.emitter.SetWindow(win)
}

// RunASV copies bundled sqlite into userDB (idempotent).
func (s *Service) RunASV(userDB, resourcesDir string) error {
	s.setStatus(progress.StepQueued, "pending", "")

	bundledSQLite := filepath.Join(resourcesDir, bundled.SQLiteFile)
	confPath := filepath.Join(resourcesDir, bundled.ConfFile)

	s.setStatus(progress.StepCreatingSchema, "creating_schema", "")
	if err := bundled.ImportASV(userDB, bundledSQLite, confPath); err != nil {
		s.setStatus(progress.StepFailed, "failed", err.Error())
		return err
	}

	s.setStatus(progress.StepIndexing, "indexing", "")
	s.setStatus(progress.StepReady, "ready", "")

	return nil
}
