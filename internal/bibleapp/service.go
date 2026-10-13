package bibleapp

import (
	"errors"
	"fmt"
	"path/filepath"

	"boanerges/internal/bible/bootstrap"
	"boanerges/internal/bible/crossrefs"
	"boanerges/internal/bible/read"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/search"
	"boanerges/internal/bible/translations"
	"boanerges/internal/dictionary"

	"github.com/egoist/mygo"
)

var errNotFound = errors.New("not found")

// BibleService exposes bible reads to the frontend.
type BibleService struct {
	dbPath    string
	bootstrap *bootstrap.Service
}

func NewBibleService(dbPath string, boot *bootstrap.Service) *BibleService {
	return &BibleService{dbPath: dbPath, bootstrap: boot}
}

func (s *BibleService) openMgr() (*schema.Manager, error) {
	return schema.Open(s.dbPath)
}

func (s *BibleService) openReader() (*read.Reader, *schema.Manager, error) {
	mgr, err := s.openMgr()
	if err != nil {
		return nil, nil, err
	}

	return read.NewReader(mgr.DB(), mgr, nil), mgr, nil
}

func (s *BibleService) requireReadyTranslation(abbrev string) error {
	mgr, err := s.openMgr()
	if err != nil {
		return err
	}
	defer mgr.Close()

	store := translations.NewStore(mgr)
	rec, err := store.FindByAbbrev(abbrev)
	if err != nil {
		return err
	}
	if rec == nil || !rec.IsReady() {
		return fmt.Errorf("%w: translation %q is not installed", errNotFound, abbrev)
	}

	return nil
}

type booksResponse struct {
	Books []read.Book `json:"books"`
}

type chapterResponse struct {
	Chapter read.Chapter `json:"chapter"`
}

type crossRefsResponse struct {
	References []crossrefs.Reference `json:"references"`
}

// Books lists books for an installed translation.
func (s *BibleService) Books(translation string) (booksResponse, error) {
	if err := s.requireReadyTranslation(translation); err != nil {
		return booksResponse{}, err
	}
	mgr, err := s.openMgr()
	if err != nil {
		return booksResponse{}, err
	}
	defer mgr.Close()

	books, err := read.ListBooks(mgr.DB(), mgr, translation)
	if err != nil {
		return booksResponse{}, err
	}

	return booksResponse{Books: books}, nil
}

// Chapter returns a chapter for translation/book (legacy JSON shape).
func (s *BibleService) Chapter(translation, book string, chapter int) (chapterResponse, error) {
	if err := s.requireReadyTranslation(translation); err != nil {
		return chapterResponse{}, err
	}
	reader, mgr, err := s.openReader()
	if err != nil {
		return chapterResponse{}, err
	}
	defer mgr.Close()

	ch, err := reader.Chapter(translation, book, chapter)
	if err != nil {
		return chapterResponse{}, err
	}

	return chapterResponse{Chapter: ch}, nil
}

// Search runs FTS query across ready translations.
func (s *BibleService) Search(query string, translation string, limit, offset int) (search.Response, error) {
	mgr, err := s.openMgr()
	if err != nil {
		return search.Response{}, err
	}
	defer mgr.Close()

	var trans *string
	if translation != "" {
		trans = &translation
	}

	return search.Search(mgr.DB(), mgr, translations.NewStore(mgr), query, trans, limit, offset)
}

// CrossReferences returns ranked cross references for a verse.
func (s *BibleService) CrossReferences(book string, chapter, verse int) (crossRefsResponse, error) {
	mgr, err := s.openMgr()
	if err != nil {
		return crossRefsResponse{}, err
	}
	defer mgr.Close()

	svc := crossrefs.New(mgr.DB())
	if !svc.IsImported() {
		return crossRefsResponse{}, fmt.Errorf("cross-references bundle is not imported")
	}
	refs, err := svc.ForVerse(book, chapter, verse)
	if err != nil {
		return crossRefsResponse{}, err
	}

	return crossRefsResponse{References: refs}, nil
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

// DictionaryService exposes Webster 1828 lookup/suggest.
type DictionaryService struct {
	dbPath string
}

func NewDictionaryService(dbPath string) *DictionaryService {
	return &DictionaryService{dbPath: dbPath}
}

type suggestResponse struct {
	Suggestions []string `json:"suggestions"`
}

type importingResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Suggest returns dictionary prefix matches.
func (s *DictionaryService) Suggest(query string, limit int) (any, error) {
	svc, err := s.openDict()
	if err != nil {
		return nil, err
	}
	if !svc.IsImported() {
		return importingResponse{Status: "importing", Message: "Dictionary is still being set up."}, errImporting
	}

	return suggestResponse{Suggestions: svc.Suggest(query, limit)}, nil
}

// Lookup returns dictionary entry for a word.
func (s *DictionaryService) Lookup(word string) (any, error) {
	svc, err := s.openDict()
	if err != nil {
		return nil, err
	}
	if !svc.IsImported() {
		return importingResponse{Status: "importing", Message: "Dictionary is still being set up."}, errImporting
	}

	return svc.Lookup(word), nil
}

var errImporting = errors.New("importing")

func (s *DictionaryService) openDict() (*dictionary.Service, error) {
	mgr, err := schema.Open(s.dbPath)
	if err != nil {
		return nil, err
	}
	defer mgr.Close()

	return dictionary.New(mgr.DB()), nil
}
