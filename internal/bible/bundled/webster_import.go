package bundled

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"
)

var websterImportMu sync.Mutex

func logWebster(format string, args ...any) {
	log.Printf("webster: "+format, args...)
}

func dictionaryEntryCount(db *sql.DB) (int, error) {
	var n int
	err := db.QueryRow(`SELECT COUNT(*) FROM dictionary_entries`).Scan(&n)
	return n, err
}

func commitWebsterImportMeta(db *sql.DB) error {
	n, err := dictionaryEntryCount(db)
	if err != nil {
		return err
	}
	if n == 0 {
		clearWebsterImportMeta(db)
		return fmt.Errorf("webster dictionary: import produced 0 entries")
	}
	completedAt := time.Now().UTC().Format(time.RFC3339)
	_, err = db.Exec(`INSERT INTO import_meta (key, completed_at) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET completed_at = excluded.completed_at`,
		websterMetaKey, completedAt)
	if err != nil {
		return err
	}
	logWebster("installed %d entries", n)
	return nil
}
