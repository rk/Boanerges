package schema

import "fmt"

// IndexBookCoverage sets has_content from 1:1 verse text and tightens chapters to max imported chapter.
func (m *Manager) IndexBookCoverage(abbrev string) error {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}
	if err := m.EnsureBooksHasContentColumn(a); err != nil {
		return err
	}
	books, err := m.BooksTable(a)
	if err != nil {
		return err
	}
	verses, err := m.VersesTable(a)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(fmt.Sprintf(`
		UPDATE %s SET
			has_content = EXISTS(
				SELECT 1 FROM %s v
				WHERE v.book_id = %s.id AND v.chapter = 1 AND v.verse = 1 AND TRIM(v.text) <> ''
			),
			chapters = COALESCE(
				(SELECT MAX(v.chapter) FROM %s v WHERE v.book_id = %s.id),
				chapters
			)`, books, verses, books, verses, books))

	return err
}

// EnsureBooksHasContentColumn adds has_content to legacy books tables.
func (m *Manager) EnsureBooksHasContentColumn(abbrev string) error {
	books, err := m.BooksTable(abbrev)
	if err != nil {
		return err
	}
	rows, err := m.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", books))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var cid, notnull, pk int
		var name, ctype string
		var dflt interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
			return err
		}
		if name == "has_content" {
			return rows.Err()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = m.db.Exec(fmt.Sprintf(
		`ALTER TABLE %s ADD COLUMN has_content INTEGER NOT NULL DEFAULT 1`, books))

	return err
}
