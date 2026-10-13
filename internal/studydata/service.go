package studydata

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/egoist/mygo"
)

var bookIDPattern = regexp.MustCompile(`^[a-z0-9]{1,10}$`)

// Service stores notes, scribe, and verse lists under user data (legacy paths).
type Service struct {
	root string
}

func NewService() (*Service, error) {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return nil, err
	}

	return &Service{root: dir}, nil
}

// NewServiceForTest uses an explicit root (legacy-compatible on-disk layout).
func NewServiceForTest(root string) *Service {
	return &Service{root: root}
}

func (s *Service) NotesGet(bookID string, chapter int) (string, error) {
	if !bookIDPattern.MatchString(bookID) {
		return "", fmt.Errorf("invalid book id")
	}
	path := filepath.Join(s.root, "notes", bookID, fmt.Sprintf("%d.txt", chapter))
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}

	return string(raw), err
}

func (s *Service) NotesPut(bookID string, chapter int, content string) error {
	if !bookIDPattern.MatchString(bookID) {
		return fmt.Errorf("invalid book id")
	}
	path := filepath.Join(s.root, "notes", bookID, fmt.Sprintf("%d.txt", chapter))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0o644)
}

type ScribeVerse struct {
	Verse           int    `json:"verse"`
	Text            string `json:"text"`
	ParagraphStart  bool   `json:"paragraphStart,omitempty"`
}

func (s *Service) ScribeGet(bookID string, chapter int) ([]ScribeVerse, error) {
	if !bookIDPattern.MatchString(bookID) {
		return nil, fmt.Errorf("invalid book id")
	}
	path := filepath.Join(s.root, "scribe", bookID, fmt.Sprintf("%d.json", chapter))
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []ScribeVerse{}, nil
	}
	if err != nil {
		return nil, err
	}
	var verses []ScribeVerse
	if err := json.Unmarshal(raw, &verses); err != nil {
		return []ScribeVerse{}, nil
	}

	return verses, nil
}

func (s *Service) ScribePut(bookID string, chapter int, verses []ScribeVerse) error {
	if !bookIDPattern.MatchString(bookID) {
		return fmt.Errorf("invalid book id")
	}
	path := filepath.Join(s.root, "scribe", bookID, fmt.Sprintf("%d.json", chapter))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.Marshal(verses)
	if err != nil {
		return err
	}

	return os.WriteFile(path, raw, 0o644)
}

type VerseListEntry struct {
	BookID  string `json:"bookId"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
}

type VerseListSummary struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	UpdatedAt string `json:"updatedAt"`
}

type VerseList struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Entries   []VerseListEntry `json:"entries"`
	UpdatedAt string           `json:"updatedAt"`
}

func (s *Service) VerseListIndex() ([]VerseListSummary, error) {
	index, err := s.readVerseListIndex()
	if err != nil {
		return nil, err
	}
	sort.Slice(index, func(i, j int) bool { return index[i].UpdatedAt > index[j].UpdatedAt })

	return index, nil
}

func (s *Service) VerseListGet(id string) (VerseList, error) {
	path := filepath.Join(s.root, "verse-lists", id+".json")
	raw, err := os.ReadFile(path)
	if err != nil {
		return VerseList{}, fmt.Errorf("verse list not found")
	}
	var data map[string]any
	if err := json.Unmarshal(raw, &data); err != nil {
		return VerseList{}, err
	}

	return normalizeList(id, data), nil
}

func (s *Service) VerseListSave(title string, entries []VerseListEntry, id string) (VerseList, error) {
	if strings.TrimSpace(title) == "" {
		return VerseList{}, fmt.Errorf("title required")
	}
	for _, entry := range entries {
		if !bookIDPattern.MatchString(entry.BookID) {
			return VerseList{}, fmt.Errorf("invalid book id")
		}
		if entry.Chapter <= 0 || entry.Verse <= 0 {
			return VerseList{}, fmt.Errorf("invalid reference")
		}
	}
	if id == "" {
		id = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	list := VerseList{
		ID: id, Title: strings.TrimSpace(title), Entries: entries,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	path := filepath.Join(s.root, "verse-lists", id+".json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return VerseList{}, err
	}
	raw, err := json.Marshal(list)
	if err != nil {
		return VerseList{}, err
	}
	if err := os.WriteFile(path, raw, 0o644); err != nil {
		return VerseList{}, err
	}
	index, _ := s.readVerseListIndex()
	found := false
	for i, item := range index {
		if item.ID == id {
			index[i].Title = list.Title
			index[i].UpdatedAt = list.UpdatedAt
			found = true
			break
		}
	}
	if !found {
		index = append(index, VerseListSummary{ID: id, Title: list.Title, UpdatedAt: list.UpdatedAt})
	}
	if err := s.writeVerseListIndex(index); err != nil {
		return VerseList{}, err
	}

	return list, nil
}

func (s *Service) VerseListDelete(id string) error {
	_ = os.Remove(filepath.Join(s.root, "verse-lists", id+".json"))
	index, err := s.readVerseListIndex()
	if err != nil {
		return err
	}
	filtered := index[:0]
	for _, item := range index {
		if item.ID != id {
			filtered = append(filtered, item)
		}
	}

	return s.writeVerseListIndex(filtered)
}

func (s *Service) readVerseListIndex() ([]VerseListSummary, error) {
	path := filepath.Join(s.root, "verse-lists", "index.json")
	raw, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []VerseListSummary{}, nil
	}
	if err != nil {
		return nil, err
	}
	var index []VerseListSummary
	if err := json.Unmarshal(raw, &index); err != nil {
		return []VerseListSummary{}, nil
	}

	return index, nil
}

func (s *Service) writeVerseListIndex(index []VerseListSummary) error {
	path := filepath.Join(s.root, "verse-lists", "index.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.Marshal(index)
	if err != nil {
		return err
	}

	return os.WriteFile(path, raw, 0o644)
}

func normalizeList(id string, data map[string]any) VerseList {
	title, _ := data["title"].(string)
	updated, _ := data["updatedAt"].(string)
	var entries []VerseListEntry
	if rawEntries, ok := data["entries"].([]any); ok {
		for _, item := range rawEntries {
			m, _ := item.(map[string]any)
			entries = append(entries, VerseListEntry{
				BookID:  fmt.Sprint(m["bookId"]),
				Chapter: intFromAny(m["chapter"]),
				Verse:   intFromAny(m["verse"]),
			})
		}
	}

	return VerseList{ID: id, Title: title, Entries: entries, UpdatedAt: updated}
}

func intFromAny(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	default:
		return 0
	}
}
