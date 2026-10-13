package schema

import (
	"fmt"

	"boanerges/internal/bible/markup"
)

const plainTextBackfillChunk = 500

// BackfillPlainTextIfEmpty fills plain_text from markup in chunks (separate from FTS rebuild).
func (m *Manager) BackfillPlainTextIfEmpty(abbrev string, formatter *markup.Formatter) error {
	if formatter == nil {
		formatter = markup.DefaultFormatter()
	}
	verses, err := m.VersesTable(abbrev)
	if err != nil {
		return err
	}

	var empty int
	if err := m.db.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE plain_text IS NULL OR plain_text = ''`, verses,
	)).Scan(&empty); err != nil {
		return err
	}
	if empty == 0 {
		return nil
	}

	for {
		rows, err := m.db.Query(fmt.Sprintf(
			`SELECT id, text FROM %s WHERE plain_text IS NULL OR plain_text = '' ORDER BY id LIMIT ?`, verses,
		), plainTextBackfillChunk)
		if err != nil {
			return err
		}

		type row struct {
			id   int64
			text string
		}
		var batch []row
		for rows.Next() {
			var r row
			if err := rows.Scan(&r.id, &r.text); err != nil {
				rows.Close()
				return err
			}
			batch = append(batch, r)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return err
		}
		if len(batch) == 0 {
			return nil
		}

		tx, err := m.db.Begin()
		if err != nil {
			return err
		}
		stmt, err := tx.Prepare(fmt.Sprintf("UPDATE %s SET plain_text=? WHERE id=?", verses))
		if err != nil {
			_ = tx.Rollback()
			return err
		}
		for _, r := range batch {
			plain := formatter.ToPlainText(r.text)
			if plain == "" {
				plain = r.text
			}
			if _, err := stmt.Exec(plain, r.id); err != nil {
				_ = stmt.Close()
				_ = tx.Rollback()
				return err
			}
		}
		_ = stmt.Close()
		if err := tx.Commit(); err != nil {
			return err
		}
	}
}
