package studydata

import "sync"

// StudyData exposes study file stores to the frontend (lazy-init for mygo generate).
type StudyData struct {
	once sync.Once
	svc  *Service
	err  error
}

func (b *StudyData) open() (*Service, error) {
	b.once.Do(func() {
		b.svc, b.err = NewService()
	})

	return b.svc, b.err
}

func (b *StudyData) NotesGet(bookID string, chapter int) (string, error) {
	s, err := b.open()
	if err != nil {
		return "", err
	}

	return s.NotesGet(bookID, chapter)
}

func (b *StudyData) NotesPut(bookID string, chapter int, content string) error {
	s, err := b.open()
	if err != nil {
		return err
	}

	return s.NotesPut(bookID, chapter, content)
}

func (b *StudyData) ScribeGet(bookID string, chapter int) ([]ScribeVerse, error) {
	s, err := b.open()
	if err != nil {
		return nil, err
	}

	return s.ScribeGet(bookID, chapter)
}

func (b *StudyData) ScribePut(bookID string, chapter int, verses []ScribeVerse) error {
	s, err := b.open()
	if err != nil {
		return err
	}

	return s.ScribePut(bookID, chapter, verses)
}

func (b *StudyData) VerseListIndex() ([]VerseListSummary, error) {
	s, err := b.open()
	if err != nil {
		return nil, err
	}

	return s.VerseListIndex()
}

func (b *StudyData) VerseListGet(id string) (VerseList, error) {
	s, err := b.open()
	if err != nil {
		return VerseList{}, err
	}

	return s.VerseListGet(id)
}

func (b *StudyData) VerseListSave(title string, entries []VerseListEntry, id string) (VerseList, error) {
	s, err := b.open()
	if err != nil {
		return VerseList{}, err
	}

	return s.VerseListSave(title, entries, id)
}

func (b *StudyData) VerseListDelete(id string) error {
	s, err := b.open()
	if err != nil {
		return err
	}

	return s.VerseListDelete(id)
}
