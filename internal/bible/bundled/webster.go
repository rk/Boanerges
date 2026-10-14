package bundled

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"boanerges/internal/bible/schema"
)

const websterMetaKey = "webster1828_dictionary_v1"

// ImportWebster1828 loads bundled dictionary JSON (idempotent).
func ImportWebster1828(dbPath, jsonPath string) error {
	if _, err := os.Stat(jsonPath); err != nil {
		return fmt.Errorf("webster dictionary: %w", err)
	}

	mgr, err := schema.Open(dbPath)
	if err != nil {
		return err
	}
	defer mgr.Close()

	if err := mgr.CreateAppSchema(); err != nil {
		return err
	}

	if WebsterInstalled(mgr.DB()) {
		return nil
	}
	clearWebsterImportMeta(mgr.DB())

	raw, err := os.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var entries []map[string]any
	if err := json.Unmarshal(raw, &entries); err != nil {
		return fmt.Errorf("decode webster json: %w", err)
	}

	db := mgr.DB()
	if _, err := db.Exec(`DELETE FROM dictionary_entries`); err != nil {
		return err
	}

	batch := make([]dictRow, 0, 1000)
	for _, entry := range entries {
		word, _ := entry["word"].(string)
		word = trim(word)
		if word == "" {
			continue
		}
		defs, ok := entry["definitions"].([]any)
		if !ok {
			continue
		}
		defStrings := make([]string, 0, len(defs))
		for _, d := range defs {
			defStrings = append(defStrings, fmt.Sprint(d))
		}
		encoded, err := json.Marshal(defStrings)
		if err != nil {
			return err
		}
		pos := ""
		if p, ok := entry["pos"].(string); ok {
			pos = p
		}
		batch = append(batch, dictRow{
			Word: word, WordKey: strings.ToLower(word), POS: nullIfEmpty(pos), Definitions: string(encoded),
		})
		if len(batch) >= 1000 {
			if err := insertDictBatch(db, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := insertDictBatch(db, batch); err != nil {
			return err
		}
	}

	return commitWebsterImportMeta(db)
}

type dictRow struct {
	Word        string
	WordKey     string
	POS         any
	Definitions string
}

func insertDictBatch(db *sql.DB, batch []dictRow) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO dictionary_entries (word, word_key, part_of_speech, definitions) VALUES (?, ?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, row := range batch {
		if _, err := stmt.Exec(row.Word, row.WordKey, row.POS, row.Definitions); err != nil {
			_ = stmt.Close()
			_ = tx.Rollback()
			return err
		}
	}
	_ = stmt.Close()

	return tx.Commit()
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}

	return s
}
