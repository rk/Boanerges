package bundled

import (
	"database/sql"

	"boanerges/internal/bible/schema"
)

// WebsterInstalled is true when import_meta is set and dictionary_entries has data.
func WebsterInstalled(db *sql.DB) bool {
	var completed sql.NullString
	if err := db.QueryRow(`SELECT completed_at FROM import_meta WHERE key = ?`, websterMetaKey).Scan(&completed); err != nil || !completed.Valid {
		return false
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM dictionary_entries`).Scan(&n); err != nil || n == 0 {
		return false
	}
	return true
}

func clearWebsterImportMeta(db *sql.DB) {
	_, _ = db.Exec(`DELETE FROM import_meta WHERE key = ?`, websterMetaKey)
}

// ResetWebster1828 clears Webster import state so EnsureWebster1828 can run again.
func ResetWebster1828(userDB string) error {
	mgr, err := schema.Open(userDB)
	if err != nil {
		return err
	}
	defer mgr.Close()
	if err := mgr.CreateAppSchema(); err != nil {
		return err
	}
	clearWebsterImportMeta(mgr.DB())
	_, err = mgr.DB().Exec(`DELETE FROM dictionary_entries`)
	return err
}
