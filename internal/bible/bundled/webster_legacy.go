package bundled

import (
	"database/sql"
	"fmt"
	"time"

	"boanerges/internal/bible/schema"
)

// ImportWebster1828FromLegacyDB copies dictionary_entries from the reference NativePHP sqlite.
func ImportWebster1828FromLegacyDB(userDB, legacyDB string) error {
	mgr, err := schema.Open(userDB)
	if err != nil {
		return err
	}
	defer mgr.Close()

	if err := mgr.CreateAppSchema(); err != nil {
		return err
	}

	var completed sql.NullString
	if err := mgr.DB().QueryRow(`SELECT completed_at FROM import_meta WHERE key = ?`, websterMetaKey).Scan(&completed); err == nil && completed.Valid {
		return nil
	}

	db := mgr.DB()
	if _, err := db.Exec(`ATTACH DATABASE ? AS webster_src`, legacyDB); err != nil {
		return fmt.Errorf("attach legacy dictionary db: %w", err)
	}
	defer func() { _, _ = db.Exec(`DETACH webster_src`) }()

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM webster_src.dictionary_entries`).Scan(&count); err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("legacy dictionary has no entries")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM dictionary_entries`); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`
		INSERT INTO dictionary_entries (word, word_key, part_of_speech, definitions)
		SELECT word, word_key, part_of_speech, definitions FROM webster_src.dictionary_entries`); err != nil {
		_ = tx.Rollback()
		return err
	}

	completedAt := time.Now().UTC().Format(time.RFC3339)
	var legacyCompleted sql.NullString
	_ = db.QueryRow(`SELECT completed_at FROM webster_src.import_meta WHERE key = ?`, websterMetaKey).Scan(&legacyCompleted)
	if legacyCompleted.Valid {
		completedAt = legacyCompleted.String
	}

	if _, err := tx.Exec(`INSERT INTO import_meta (key, completed_at) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET completed_at = excluded.completed_at`,
		websterMetaKey, completedAt); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// EnsureWebster1828 loads Webster when missing: bundled JSON, then reference legacy DB.
func EnsureWebster1828(userDB, resourcesDir string) error {
	mgr, err := schema.Open(userDB)
	if err != nil {
		return err
	}
	var done sql.NullString
	_ = mgr.DB().QueryRow(`SELECT completed_at FROM import_meta WHERE key = ?`, websterMetaKey).Scan(&done)
	_ = mgr.Close()
	if done.Valid {
		return nil
	}

	if jsonPath, err := FindWebsterDictionaryJSON(resourcesDir); err == nil {
		return ImportWebster1828(userDB, jsonPath)
	}
	if legacyDB, err := FindLegacyDictionaryDB(); err == nil {
		return ImportWebster1828FromLegacyDB(userDB, legacyDB)
	}

	return fmt.Errorf("webster dictionary: no bundled JSON or legacy database found")
}
