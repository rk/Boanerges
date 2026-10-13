package appchrome

import "github.com/egoist/mygo"

// StudySettingsChanged is emitted on the main window when the flyout picker saves location.
var StudySettingsChanged = mygo.NewEvent[struct{}]("studySettingsChanged")

// OpenBookChapterFlyout shows the book/chapter picker in a native flyout (not clipped by the webview).
// Anchor is in DIPs relative to the main window content area, e.g. from getBoundingClientRect().
func (s *AppChrome) OpenBookChapterFlyout(anchorX, anchorY, anchorWidth, anchorHeight int) {
	s.mu.Lock()
	parent := s.win
	existing := s.bookFlyout
	s.mu.Unlock()

	if parent == nil || parent.IsDestroyed() {
		return
	}

	if existing != nil && !existing.IsDestroyed() {
		s.CloseBookChapterFlyout()
		return
	}

	flyout := mygo.NewFlyout(mygo.FlyoutOptions{
		Parent:    parent,
		Anchor:    mygo.Rectangle{X: anchorX, Y: anchorY, Width: anchorWidth, Height: anchorHeight},
		URL:       "/?flyout=book-chapter",
		Width:     580,
		Height:    440,
		Gap:       4,
		Focusable: true,
		Popover:   true,
		Shadow:    true,
		Placement: mygo.PlacementBottomStart,
	})

	s.mu.Lock()
	s.bookFlyout = flyout
	s.mu.Unlock()

	flyout.OnClosed(func() {
		s.mu.Lock()
		if s.bookFlyout == flyout {
			s.bookFlyout = nil
		}
		main := s.win
		s.mu.Unlock()
		if main != nil && !main.IsDestroyed() {
			_ = StudySettingsChanged.Emit(main, struct{}{})
		}
	})
}

// CloseBookChapterFlyout dismisses the picker flyout if open.
func (s *AppChrome) CloseBookChapterFlyout() {
	s.mu.Lock()
	f := s.bookFlyout
	s.bookFlyout = nil
	s.mu.Unlock()

	if f != nil && !f.IsDestroyed() {
		f.Close()
	}
}
