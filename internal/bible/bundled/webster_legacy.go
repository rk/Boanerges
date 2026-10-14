package bundled

import (
	"fmt"

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

	if WebsterInstalled(mgr.DB()) {
		return nil
	}
	clearWebsterImportMeta(mgr.DB())

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

	// ponytail: one INSERT…SELECT autocommit; explicit tx + attached DB was stalling 100k+ rows in tests.
	_, _ = db.Exec(`PRAGMA synchronous=OFF`)
	if _, err := db.Exec(`DELETE FROM dictionary_entries`); err != nil {
		return err
	}
	if _, err := db.Exec(`
		INSERT INTO dictionary_entries (word, word_key, part_of_speech, definitions)
		SELECT word, word_key, part_of_speech, definitions FROM webster_src.dictionary_entries`); err != nil {
		return err
	}
	return commitWebsterImportMeta(db)
}

// EnsureWebster1828 loads Webster when missing: bundled JSON, then reference legacy DB.
func EnsureWebster1828(userDB, resourcesDir string) error {
	websterImportMu.Lock()
	defer websterImportMu.Unlock()

	mgr, err := schema.Open(userDB)
	if err != nil {
		return err
	}
	if WebsterInstalled(mgr.DB()) {
		_ = mgr.Close()
		return nil
	}
	clearWebsterImportMeta(mgr.DB())
	_ = mgr.Close()

	logWebster("import needed for %s", userDB)

	if jsonPath, err := FindWebsterDictionaryJSON(resourcesDir); err == nil {
		logWebster("importing from JSON %s", jsonPath)
		if err := ImportWebster1828(userDB, jsonPath); err != nil {
			logWebster("JSON import failed: %v", err)
			return err
		}
		return nil
	}
	if legacyDB, err := FindLegacyDictionaryDB(); err == nil {
		logWebster("importing from legacy DB %s", legacyDB)
		if err := ImportWebster1828FromLegacyDB(userDB, legacyDB); err != nil {
			logWebster("legacy import failed: %v", err)
			return err
		}
		return nil
	}

	err = fmt.Errorf("webster dictionary: no bundled JSON or legacy database found")
	logWebster("%v (resourcesDir=%q)", err, resourcesDir)
	return err
}
