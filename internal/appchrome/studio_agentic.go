//go:build agentic

package appchrome

import "github.com/egoist/mygo"

func (s *AppChrome) setStudioWin(win *mygo.Window) {
	s.mu.Lock()
	if prev := s.studioWin; prev != nil && !prev.IsDestroyed() && prev != win {
		prev.Close()
	}
	s.studioWin = win
	s.mu.Unlock()
}

func (s *AppChrome) clearStudioWin(win *mygo.Window) {
	s.mu.Lock()
	if s.studioWin == win {
		s.studioWin = nil
	}
	s.mu.Unlock()
}

// OpenAgenticStudio opens or focuses the Agentic Studio window.
func (s *AppChrome) OpenAgenticStudio() error {
	s.mu.Lock()
	existing := s.studioWin
	s.mu.Unlock()

	if existing != nil && !existing.IsDestroyed() {
		existing.Focus()
		return nil
	}

	studio := mygo.NewWindow(mygo.WindowOptions{
		Title:  "Boanerges Studio",
		Width:  1000,
		Height: 720,
		URL:    "/?view=studio",
		Page: mygo.PageOptions{
			DevTools: mygo.DevToolsEnabled,
		},
	})

	s.setStudioWin(studio)
	studio.OnClosed(func() {
		s.clearStudioWin(studio)
	})
	return nil
}
