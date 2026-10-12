package appchrome

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/egoist/mygo"
)

// AppChrome exposes desktop print for the main window.
type AppChrome struct {
	mu  sync.Mutex
	win *mygo.Window
}

func NewAppChrome() *AppChrome {
	return &AppChrome{}
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
