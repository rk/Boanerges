package crossrefs

import (
	"database/sql"
	"fmt"

	"boanerges/internal/bible"
	"boanerges/internal/bible/openbible"
)

const MetaKey = "cross_references_v2"

type Reference struct {
	Rank     int    `json:"rank"`
	BookID   string `json:"bookId"`
	BookName string `json:"bookName"`
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse"`
	EndVerse *int   `json:"endVerse"`
}

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) IsImported() bool {
	var completed sql.NullString
	err := s.db.QueryRow(`SELECT completed_at FROM import_meta WHERE key = ?`, MetaKey).Scan(&completed)

	return err == nil && completed.Valid
}

func (s *Service) ForVerse(bookID string, chapter, verse int) ([]Reference, error) {
	canon, ok := bible.NormalizeBookID(bookID)
	if !ok {
		return nil, nil
	}

	var sourceID int
	err := s.db.QueryRow(`
		SELECT id FROM cross_reference_verses WHERE book_id = ? AND chapter = ? AND verse = ?`,
		canon, chapter, verse).Scan(&sourceID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(`
		SELECT rank, target_start_id, target_end_id FROM cross_references
		WHERE source_verse_id = ? ORDER BY rank LIMIT 100`, sourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Reference
	for rows.Next() {
		var rank, targetStart int
		var targetEnd sql.NullInt64
		if err := rows.Scan(&rank, &targetStart, &targetEnd); err != nil {
			return nil, err
		}
		target := openbible.FromVerseID(targetStart)
		if target == nil {
			continue
		}
		var endVerse *int
		if targetEnd.Valid && int(targetEnd.Int64) != targetStart {
			if end := openbible.FromVerseID(int(targetEnd.Int64)); end != nil {
				endVerse = &end.Verse
			}
		}
		out = append(out, Reference{
			Rank: rank, BookID: target.BookID, BookName: bible.DisplayName(target.BookID),
			Chapter: target.Chapter, Verse: target.Verse, EndVerse: endVerse,
		})
	}

	return out, rows.Err()
}

func VerseID(db *sql.DB, bookID string, chapter, verse int) (int, error) {
	canon, ok := bible.NormalizeBookID(bookID)
	if !ok {
		return 0, fmt.Errorf("unknown book")
	}
	var id int
	err := db.QueryRow(`
		SELECT id FROM cross_reference_verses WHERE book_id = ? AND chapter = ? AND verse = ?`,
		canon, chapter, verse).Scan(&id)

	return id, err
}
