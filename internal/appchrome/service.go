package appchrome

import (
	"os"
	"path/filepath"
	"sync"

	"boanerges/internal/settings"
	"boanerges/internal/studyprint"

	"github.com/egoist/mygo"
)

// AppChrome exposes desktop print for the main window.
type AppChrome struct {
	mu                   sync.Mutex
	win                  *mygo.Window
	printPreviewWin      *mygo.Window
	previewStudy         settings.StudySettings
	previewReadability   settings.ReadabilitySettings
	bookFlyout           *mygo.Window
	verseFlyout          *mygo.Window
	dbPath               string
}

func NewAppChrome(dbPath string) *AppChrome {
	return &AppChrome{dbPath: dbPath}
}

func (s *AppChrome) SetWindow(win *mygo.Window) {
	s.mu.Lock()
	s.win = win
	s.mu.Unlock()
}

func (s *AppChrome) Print() {
	s.mu.Lock()
	win := s.win
	s.mu.Unlock()
	if win == nil {
		return
	}

	win.Page().Print()
}

// PrintToPDF writes a PDF export to the user downloads folder (ponytail: fixed path).
func (s *AppChrome) PrintToPDF(filename string) (string, error) {
	s.mu.Lock()
	win := s.win
	s.mu.Unlock()
	if win == nil {
		return "", nil
	}

	pdf, err := win.Page().PrintToPDF(mygo.PDFOptions{})
	if err != nil {
		return "", err
	}
	if filename == "" {
		filename = "study.pdf"
	}
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return "", err
	}
	dest := filepath.Join(dir, filename)
	if err := os.WriteFile(dest, pdf, 0o644); err != nil {
		return "", err
	}

	return dest, nil
}

// StudyPrintHTML builds printable HTML for the current study layout (primary bible + optional notes).
func (s *AppChrome) StudyPrintHTML(study settings.StudySettings, readability settings.ReadabilitySettings, includeUserWork bool) (string, error) {
	return studyprint.BuildStudyHTML(s.dbPath, study, readability, includeUserWork)
}

// WriteStudyPrintFile saves printable HTML under user data and returns the path.
func (s *AppChrome) WriteStudyPrintFile(study settings.StudySettings, readability settings.ReadabilitySettings, includeUserWork bool, filename string) (string, error) {
	html, err := studyprint.BuildStudyHTML(s.dbPath, study, readability, includeUserWork)
	if err != nil {
		return "", err
	}
	if filename == "" {
		filename = "study-print.html"
	}
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return "", err
	}
	dest := filepath.Join(dir, filename)
	if err := os.WriteFile(dest, []byte(html), 0o644); err != nil {
		return "", err
	}

	return dest, nil
}

// SyncMenuState updates native menu radio/checkbox items.
func (s *AppChrome) SyncMenuState(columnCount int, scrollSync bool) {
	SyncMenuState(columnCount, scrollSync)
}

// OpenDevTools toggles the web inspector on the focused window, or the main window.
func (s *AppChrome) OpenDevTools() {
	page := devToolsPage(s)
	if page == nil {
		return
	}
	page.ToggleDevTools()
}

func devToolsPage(s *AppChrome) *mygo.Page {
	if focused := mygo.FocusedWindow(); focused != nil {
		if p := focused.Page(); p != nil {
			return p
		}
	}
	s.mu.Lock()
	preview := s.printPreviewWin
	main := s.win
	s.mu.Unlock()
	if preview != nil && !preview.IsDestroyed() {
		if p := preview.Page(); p != nil {
			return p
		}
	}
	if main != nil && !main.IsDestroyed() {
		return main.Page()
	}
	return nil
}

// OpenFeedback opens the feedback mailto link in the system handler.
func (s *AppChrome) OpenFeedback() {
	_ = mygo.Shell.OpenExternal(feedbackMailto)
}
