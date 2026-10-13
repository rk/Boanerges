package appchrome

import (
	"fmt"
	"net/url"

	"github.com/egoist/mygo"
)

// VerseMenuPayload is dispatched on the main window after a verse context action.
type VerseMenuPayload struct {
	Action  string `json:"action"`
	BookID  string `json:"bookId"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Word    string `json:"word,omitempty"`
}

// VerseMenuAction is emitted on the main window when the verse-actions flyout picks an item.
var VerseMenuAction = mygo.NewEvent[VerseMenuPayload]("verseMenuAction")

// OpenVerseActionsFlyout shows verse context actions in a native flyout (not clipped by the webview).
func (s *AppChrome) OpenVerseActionsFlyout(anchorX, anchorY int, bookID string, chapter, verse int, word string) {
	s.mu.Lock()
	parent := s.win
	existing := s.verseFlyout
	s.mu.Unlock()

	if parent == nil || parent.IsDestroyed() {
		return
	}

	if existing != nil && !existing.IsDestroyed() {
		s.CloseVerseActionsFlyout()
	}

	q := url.Values{}
	q.Set("flyout", "verse-actions")
	q.Set("bookId", bookID)
	q.Set("chapter", fmt.Sprintf("%d", chapter))
	q.Set("verse", fmt.Sprintf("%d", verse))
	if word != "" {
		q.Set("word", word)
	}

	flyout := mygo.NewFlyout(mygo.FlyoutOptions{
		Parent:    parent,
		Anchor:    mygo.Rectangle{X: anchorX, Y: anchorY, Width: 1, Height: 1},
		URL:       "/?" + q.Encode(),
		Width:     220,
		Height:    200,
		Gap:       4,
		Focusable: true,
		Popover:   true,
		Shadow:    true,
		Placement: mygo.PlacementBottomStart,
	})

	s.mu.Lock()
	s.verseFlyout = flyout
	s.mu.Unlock()

	flyout.OnClosed(func() {
		s.mu.Lock()
		if s.verseFlyout == flyout {
			s.verseFlyout = nil
		}
		s.mu.Unlock()
	})
}

// CloseVerseActionsFlyout dismisses the verse menu flyout if open.
func (s *AppChrome) CloseVerseActionsFlyout() {
	s.mu.Lock()
	f := s.verseFlyout
	s.verseFlyout = nil
	s.mu.Unlock()

	if f != nil && !f.IsDestroyed() {
		f.Close()
	}
}

// PerformVerseAction closes the flyout and notifies the main window.
func (s *AppChrome) PerformVerseAction(action, bookID string, chapter, verse int, word string) {
	s.CloseVerseActionsFlyout()

	s.mu.Lock()
	main := s.win
	s.mu.Unlock()
	if main == nil || main.IsDestroyed() {
		return
	}

	_ = VerseMenuAction.Emit(main, VerseMenuPayload{
		Action:  action,
		BookID:  bookID,
		Chapter: chapter,
		Verse:   verse,
		Word:    word,
	})
}
