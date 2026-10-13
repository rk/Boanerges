package schema

import (
	"fmt"
	"strings"
)

// VerseRow is one row for bulk insert into a translation verses table.
type VerseRow struct {
	BookID    int64
	Chapter   int
	Verse     int
	Text      string
	PlainText string
}

// VerseCount returns rows in abbrev_verses (0 if table missing).
func (m *Manager) VerseCount(abbrev string) (int, error) {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return 0, err
	}
	table := a + "_verses"
	if ok, err := m.TableExists(table); err != nil || !ok {
		return 0, err
	}
	var n int
	err = m.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&n)

	return n, err
}

// FtsPopulated reports whether the FTS table has at least one indexed row.
func (m *Manager) FtsPopulated(abbrev string) (bool, error) {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return false, err
	}
	fts := a + "_verses_fts"
	if ok, err := m.TableExists(fts); err != nil || !ok {
		return false, err
	}
	var n int
	err = m.db.QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM %s", fts)).Scan(&n)

	return n > 0, err
}

// InsertVerses bulk-inserts rows in one statement (caller batches by chapter).
func (m *Manager) InsertVerses(abbrev string, rows []VerseRow) error {
	if len(rows) == 0 {
		return nil
	}
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}
	table, err := m.VersesTable(a)
	if err != nil {
		return err
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("INSERT INTO %s (book_id, chapter, verse, text, plain_text) VALUES ", table))
	args := make([]any, 0, len(rows)*5)
	for i, row := range rows {
		if i > 0 {
			b.WriteByte(',')
		}
		b.WriteString("(?,?,?,?,?)")
		args = append(args, row.BookID, row.Chapter, row.Verse, row.Text, row.PlainText)
	}

	_, err = m.db.Exec(b.String(), args...)

	return err
}
