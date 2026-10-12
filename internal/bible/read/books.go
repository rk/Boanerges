package read

import (
	"database/sql"
	"fmt"
	"strings"

	"boanerges/internal/bible"
	"boanerges/internal/bible/schema"
)

// Book matches legacy BookResource JSON.
type Book struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbrev       string `json:"abbrev"`
	Testament    string `json:"testament"`
	Chapters     int    `json:"chapters"`
	FirstChapter int    `json:"firstChapter"`
	LastChapter  int    `json:"lastChapter"`
}

func ListBooks(db *sql.DB, mgr *schema.Manager, abbrev string) ([]Book, error) {
	a, err := mgr.ValidateAbbrev(abbrev)
	if err != nil {
		return nil, err
	}
	table, err := mgr.BooksTable(a)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(fmt.Sprintf(
		`SELECT name, osis_id, testament, chapters FROM %s ORDER BY id`, table))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var name, osis, testament string
		var chapters int
		if err := rows.Scan(&name, &osis, &testament, &chapters); err != nil {
			return nil, err
		}
		id, _ := bible.NormalizeBookID(osis)
		if id == "" {
			id = strings.ToLower(osis)
		}
		tag := id
		if len(tag) > 3 {
			tag = tag[:3]
		}
		books = append(books, Book{
			ID: id, Name: name, Abbrev: strings.ToUpper(tag),
			Testament: testament, Chapters: chapters, FirstChapter: 1, LastChapter: chapters,
		})
	}

	return books, rows.Err()
}
