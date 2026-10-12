package bibleapp

import (
	"path/filepath"

	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/bible/read"
	"boanerges/internal/bible/schema"

	"github.com/egoist/mygo"
)

// BibleService exposes bible reads to the frontend.
type BibleService struct {
	dbPath    string
	bootstrap *bootstrap.Service
}

func NewBibleService(dbPath string, boot *bootstrap.Service) *BibleService {
	return &BibleService{dbPath: dbPath, bootstrap: boot}
}

func (s *BibleService) openReader() (*read.Reader, *schema.Manager, error) {
	mgr, err := schema.Open(s.dbPath)
	if err != nil {
		return nil, nil, err
	}

	return read.NewReader(mgr.DB(), mgr, nil), mgr, nil
}

// Chapter returns a chapter for translation/book (legacy JSON shape).
func (s *BibleService) Chapter(translation, book string, chapter int) (read.Chapter, error) {
	reader, mgr, err := s.openReader()
	if err != nil {
		return read.Chapter{}, err
	}
	defer mgr.Close()

	return reader.Chapter(translation, book, chapter)
}

// ASVInstallStatus returns bundled import progress for the welcome page.
func (s *BibleService) ASVInstallStatus() bootstrap.ProgressPayload {
	return s.bootstrap.ASVStatus()
}

// UserDBPath returns the app bible database path (for debugging).
func (s *BibleService) UserDBPath() string {
	return s.dbPath
}

func ResolveUserDBPath() (string, error) {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, "bible.sqlite"), nil
}
