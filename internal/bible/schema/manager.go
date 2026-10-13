package schema

import (
	"database/sql"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"boanerges/internal/bible/markup"

	_ "modernc.org/sqlite"
)

var abbrevPattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// Manager creates per-translation SQLite tables (books, verses, FTS5).
type Manager struct {
	db *sql.DB
}

func Open(path string) (*Manager, error) {
	dsn := sqliteDSN(path)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	return &Manager{db: db}, nil
}

func sqliteDSN(path string) string {
	if path == ":memory:" {
		return path
	}
	// ponytail: WAL + busy_timeout — many short-lived readers during background import.
	return fmt.Sprintf("file:%s?_pragma=busy_timeout(10000)&_pragma=journal_mode(WAL)", url.PathEscape(path))
}

func (m *Manager) Close() error {
	if m.db == nil {
		return nil
	}

	return m.db.Close()
}

// DB exposes the underlying connection for attach/copy operations.
func (m *Manager) DB() *sql.DB {
	return m.db
}

func (m *Manager) ValidateAbbrev(abbrev string) (string, error) {
	normalized := strings.ToLower(abbrev)
	if !abbrevPattern.MatchString(normalized) {
		return "", fmt.Errorf("invalid translation abbrev: %s", abbrev)
	}

	return normalized, nil
}

func (m *Manager) BooksTable(abbrev string) (string, error) {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return "", err
	}

	return a + "_books", nil
}

func (m *Manager) VersesTable(abbrev string) (string, error) {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return "", err
	}

	return a + "_verses", nil
}

func (m *Manager) FtsTable(abbrev string) (string, error) {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return "", err
	}

	return a + "_verses_fts", nil
}

func (m *Manager) TableExists(name string) (bool, error) {
	var count int
	err := m.db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, name).Scan(&count)

	return count > 0, err
}

func (m *Manager) CreateTables(abbrev string) error {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}
	books := a + "_books"
	verses := a + "_verses"

	_, err = m.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		osis_id TEXT NOT NULL,
		testament TEXT NOT NULL,
		chapters INTEGER NOT NULL,
		has_content INTEGER NOT NULL DEFAULT 1
	)`, books))
	if err != nil {
		return err
	}

	if err := m.EnsureBooksHasContentColumn(a); err != nil {
		return err
	}

	_, err = m.db.Exec(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		book_id INTEGER NOT NULL,
		chapter INTEGER NOT NULL,
		verse INTEGER NOT NULL,
		text TEXT NOT NULL,
		plain_text TEXT NOT NULL DEFAULT '',
		FOREIGN KEY (book_id) REFERENCES %s(id)
	)`, verses, books))
	if err != nil {
		return err
	}

	_, err = m.db.Exec(fmt.Sprintf(
		`CREATE INDEX IF NOT EXISTS %s_verses_lookup ON %s (book_id, chapter, verse)`,
		a, verses,
	))

	return err
}

func (m *Manager) CreateFtsTable(abbrev string) error {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}
	verses := a + "_verses"
	fts := a + "_verses_fts"

	_, err = m.db.Exec(fmt.Sprintf(`CREATE VIRTUAL TABLE IF NOT EXISTS %s USING fts5(
		plain_text,
		book_id UNINDEXED,
		chapter UNINDEXED,
		verse UNINDEXED,
		content='%s',
		content_rowid='id'
	)`, fts, verses))

	return err
}

func (m *Manager) RebuildFtsIndex(abbrev string, formatter *markup.Formatter) error {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}
	fts := a + "_verses_fts"

	// ponytail: skip row-by-row plain_text backfill on FTS build — use COALESCE below; import sets plain_text.

	if exists, err := m.TableExists(fts); err != nil {
		return err
	} else if exists {
		if _, err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", fts)); err != nil {
			return err
		}
	}

	if err := m.CreateFtsTable(a); err != nil {
		return err
	}

	verses := a + "_verses"
	_, err = m.db.Exec(fmt.Sprintf(`
		INSERT INTO %s (rowid, plain_text, book_id, chapter, verse)
		SELECT id, COALESCE(NULLIF(plain_text, ''), text), book_id, chapter, verse FROM %s`,
		fts, verses))

	return err
}

func (m *Manager) DropTables(abbrev string) error {
	a, err := m.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}

	for _, table := range []string{a + "_verses_fts", a + "_verses", a + "_books"} {
		if exists, err := m.TableExists(table); err != nil {
			return err
		} else if exists {
			if _, err := m.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
				return err
			}
		}
	}

	return nil
}

// SupportsFTS5 reports whether this SQLite build can create fts5 virtual tables.
func SupportsFTS5(db *sql.DB) bool {
	if _, err := db.Exec(`CREATE VIRTUAL TABLE IF NOT EXISTS __fts_probe USING fts5(content)`); err != nil {
		return false
	}
	_, _ = db.Exec(`DROP TABLE IF EXISTS __fts_probe`)

	return true
}
