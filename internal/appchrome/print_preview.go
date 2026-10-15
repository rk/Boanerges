package appchrome

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"boanerges/internal/settings"
	"boanerges/internal/studyprint"

	"github.com/egoist/mygo"
)

func (s *AppChrome) setPrintPreviewWin(win *mygo.Window) {
	s.mu.Lock()
	if prev := s.printPreviewWin; prev != nil && !prev.IsDestroyed() && prev != win {
		prev.Close()
	}
	s.printPreviewWin = win
	s.mu.Unlock()
}

func (s *AppChrome) clearPrintPreviewWin(win *mygo.Window) {
	s.mu.Lock()
	if s.printPreviewWin == win {
		s.printPreviewWin = nil
	}
	s.mu.Unlock()
}

// StudyPrintPreviewRefresh tells the preview page to reload after OpenStudyPrintPreview reuses the window.
var StudyPrintPreviewRefresh = mygo.NewEvent[struct{}]("studyPrintPreviewRefresh")

// PrintPreviewDocument is printable CSS and body HTML for the preview page.
type PrintPreviewDocument = studyprint.PrintDocument

// OpenStudyPrintPreview opens (or focuses) the study print preview window.
func (s *AppChrome) OpenStudyPrintPreview(study settings.StudySettings, readability settings.ReadabilitySettings) error {
	s.mu.Lock()
	s.previewStudy = study
	s.previewReadability = readability
	existing := s.printPreviewWin
	s.mu.Unlock()

	if existing != nil && !existing.IsDestroyed() {
		existing.Focus()
		_ = StudyPrintPreviewRefresh.Emit(existing, struct{}{})
		return nil
	}

	printWin := mygo.NewWindow(mygo.WindowOptions{
		Title:  "Study Print Preview",
		Width:  920,
		Height: 720,
		URL:    "/?view=print-preview",
		Page: mygo.PageOptions{
			DevTools: mygo.DevToolsEnabled,
		},
	})

	s.setPrintPreviewWin(printWin)

	printWin.OnClosed(func() {
		s.clearPrintPreviewWin(printWin)
		s.mu.Lock()
		s.previewStudy = settings.StudySettings{}
		s.mu.Unlock()
	})

	return nil
}

// PrintPreviewDocument returns printable content for the current preview session.
func (s *AppChrome) PrintPreviewDocument(includeUserWork bool) (PrintPreviewDocument, error) {
	study, readability, err := s.previewSession()
	if err != nil {
		return PrintPreviewDocument{}, err
	}
	return studyprint.BuildStudyPrintDocument(s.dbPath, study, readability, includeUserWork)
}

// PrintPreviewPrint opens the system print dialog for the preview window.
func (s *AppChrome) PrintPreviewPrint(ctx context.Context) error {
	win := mygo.CallerWindow(ctx)
	if !s.isPrintPreviewWindow(win) {
		return errors.New("print preview: call from the preview window")
	}
	win.Page().Print()
	return nil
}

// PrintPreviewSavePDF renders the preview page to PDF and returns the saved path.
func (s *AppChrome) PrintPreviewSavePDF(ctx context.Context) (string, error) {
	win := mygo.CallerWindow(ctx)
	if !s.isPrintPreviewWindow(win) {
		return "", errors.New("print preview: call from the preview window")
	}

	study, _, err := s.previewSession()
	if err != nil {
		return "", err
	}

	const marginIn = 0.5
	pdf, err := win.Page().PrintToPDF(mygo.PDFOptions{
		PageSize:   mygo.PageA4,
		Landscape:  study.ColumnCount > 1,
		Background: true,
		Margins:    &mygo.Margins{Top: marginIn, Right: marginIn, Bottom: marginIn, Left: marginIn},
	})
	if err != nil {
		return "", err
	}

	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return "", err
	}
	name := fmt.Sprintf("study-print-%d.pdf", time.Now().Unix())
	dest := filepath.Join(dir, name)
	if err := os.WriteFile(dest, pdf, 0o644); err != nil {
		return "", err
	}
	return dest, nil
}

func (s *AppChrome) previewSession() (settings.StudySettings, settings.ReadabilitySettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.previewStudy.BookID == "" {
		return settings.StudySettings{}, settings.ReadabilitySettings{}, errors.New("print preview: no active session")
	}
	return s.previewStudy, s.previewReadability, nil
}

func (s *AppChrome) isPrintPreviewWindow(win *mygo.Window) bool {
	if win == nil || win.IsDestroyed() {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.printPreviewWin != nil && s.printPreviewWin.ID() == win.ID()
}
