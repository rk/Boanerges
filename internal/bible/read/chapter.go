package read

import (
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"boanerges/internal/bible"
	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
)

var paragraphStart = regexp.MustCompile(`(?i)<(p|div)\b`)

// Verse is one verse in a chapter response.
type Verse struct {
	Number          int    `json:"number"`
	Text            string `json:"text"`
	ParagraphStart  bool   `json:"paragraphStart,omitempty"`
}

// Chapter matches legacy ChapterResource JSON.
type Chapter struct {
	Book       string  `json:"book"`
	BookAbbrev string  `json:"bookAbbrev"`
	Chapter    int     `json:"chapter"`
	Verses     []Verse `json:"verses"`
}

type Reader struct {
	db        *sql.DB
	schema    *schema.Manager
	formatter *markup.Formatter
}

func NewReader(db *sql.DB, mgr *schema.Manager, formatter *markup.Formatter) *Reader {
	if formatter == nil {
		formatter = markup.DefaultFormatter()
	}

	return &Reader{db: db, schema: mgr, formatter: formatter}
}

func (r *Reader) Chapter(abbrev, bookID string, chapter int) (Chapter, error) {
	a, err := r.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return Chapter{}, err
	}

	canon, ok := bible.NormalizeBookID(bookID)
	if !ok {
		return Chapter{}, fmt.Errorf("book not found: %s", bookID)
	}

	booksTable, err := r.schema.BooksTable(a)
	if err != nil {
		return Chapter{}, err
	}
	versesTable, err := r.schema.VersesTable(a)
	if err != nil {
		return Chapter{}, err
	}

	var book struct {
		ID       int64
		Name     string
		OsisID   string
		Chapters int
	}
	err = r.db.QueryRow(fmt.Sprintf(
		`SELECT id, name, osis_id, chapters FROM %s WHERE osis_id = ?`, booksTable,
	), canon).Scan(&book.ID, &book.Name, &book.OsisID, &book.Chapters)
	if err == sql.ErrNoRows {
		return Chapter{}, fmt.Errorf("book not found: %s", bookID)
	}
	if err != nil {
		return Chapter{}, err
	}
	if chapter < 1 || chapter > book.Chapters {
		return Chapter{}, fmt.Errorf("chapter %d out of range for %s", chapter, book.Name)
	}

	rows, err := r.db.Query(fmt.Sprintf(
		`SELECT verse, text FROM %s WHERE book_id = ? AND chapter = ? ORDER BY verse`,
		versesTable,
	), book.ID, chapter)
	if err != nil {
		return Chapter{}, err
	}
	defer rows.Close()

	var verses []Verse
	for rows.Next() {
		var num int
		var raw string
		if err := rows.Scan(&num, &raw); err != nil {
			return Chapter{}, err
		}
		v := Verse{
			Number: num,
			Text:   strings.TrimSpace(r.formatter.Format(raw)),
		}
		trimmed := strings.TrimLeft(raw, " \t\r\n")
		if paragraphStart.MatchString(trimmed) || strings.Contains(trimmed, `type="x-p"`) {
			v.ParagraphStart = true
		}
		verses = append(verses, v)
	}

	return Chapter{
		Book:       book.Name,
		BookAbbrev: strings.ToUpper(canon[:min(3, len(canon))]),
		Chapter:    chapter,
		Verses:     verses,
	}, rows.Err()
}
