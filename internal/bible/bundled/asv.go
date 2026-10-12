package bundled

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"boanerges/internal/bible/schema"
	"boanerges/sword"
)

const (
	// SQLiteFile is the bundled translation tables file (under PathResources).
	SQLiteFile = "bundled/asv.sqlite"
	ConfFile   = "sword/mods.d/asv.conf"
)

// ImportASV copies pre-built ASV tables from bundledSQLite into dbPath and registers metadata.
func ImportASV(dbPath, bundledSQLite, confPath string) error {
	if _, err := os.Stat(bundledSQLite); err != nil {
		return fmt.Errorf("bundled asv sqlite: %w", err)
	}

	mgr, err := schema.Open(dbPath)
	if err != nil {
		return err
	}
	defer mgr.Close()

	if err := mgr.CreateAppSchema(); err != nil {
		return err
	}

	exists, err := mgr.TableExists("asv_verses")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	if err := mgr.CreateTables("asv"); err != nil {
		return err
	}

	db := mgr.DB()
	attach := filepath.ToSlash(bundledSQLite)
	if _, err := db.Exec(`ATTACH DATABASE ? AS bundled`, attach); err != nil {
		return err
	}
	defer db.Exec(`DETACH DATABASE bundled`)

	for _, table := range []string{"asv_books", "asv_verses"} {
		if _, err := db.Exec(fmt.Sprintf(
			`INSERT INTO main.%s SELECT * FROM bundled.%s`, table, table,
		)); err != nil {
			return err
		}
	}

	if err := mgr.CreateFtsTable("asv"); err != nil {
		return err
	}
	if _, err := db.Exec(`INSERT INTO main.asv_verses_fts(asv_verses_fts) VALUES('rebuild')`); err != nil {
		return err
	}

	return insertTranslationRow(db, confPath)
}

func insertTranslationRow(db *sql.DB, confPath string) error {
	meta := sword.Metadata{Name: "American Standard Version (1901)"}
	if _, m, err := sword.ParseConfFile(confPath); err == nil {
		meta = m
	}

	_, err := db.Exec(`
		INSERT OR IGNORE INTO translations (
			abbrev, name, format, versification, about, version_string, version_date,
			copyright, copyright_contact, source, install_status, install_step, install_error, bundled,
			created_at, updated_at
		) VALUES (
			'asv', ?, ?, ?, ?, ?, ?,
			?, ?, ?, 'ready', 'ready', NULL, 1,
			datetime('now'), datetime('now')
		)`,
		meta.Name, meta.Format, meta.Versification, meta.About, meta.VersionString, meta.VersionDate,
		meta.Copyright, meta.CopyrightContact, meta.Source,
	)

	return err
}
