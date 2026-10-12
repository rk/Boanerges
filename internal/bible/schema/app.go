package schema

// CreateAppSchema creates shared catalog tables (translations, import_meta).
func (m *Manager) CreateAppSchema() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS translations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		abbrev TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		format TEXT,
		versification TEXT,
		about TEXT,
		version_string TEXT,
		version_date TEXT,
		copyright TEXT,
		copyright_contact TEXT,
		source TEXT,
		install_status TEXT NOT NULL DEFAULT 'pending',
		install_step TEXT,
		install_error TEXT,
		bundled INTEGER NOT NULL DEFAULT 0,
		created_at TEXT NOT NULL DEFAULT (datetime('now')),
		updated_at TEXT NOT NULL DEFAULT (datetime('now'))
	)`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS import_meta (
		key TEXT PRIMARY KEY,
		completed_at TEXT
	)`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS cross_reference_verses (
		id INTEGER PRIMARY KEY,
		book_id TEXT NOT NULL,
		chapter INTEGER NOT NULL,
		verse INTEGER NOT NULL
	)`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS cross_references (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source_verse_id INTEGER NOT NULL,
		rank INTEGER NOT NULL,
		target_start_id INTEGER NOT NULL,
		target_end_id INTEGER,
		FOREIGN KEY (source_verse_id) REFERENCES cross_reference_verses(id)
	)`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE INDEX IF NOT EXISTS cross_references_source_rank ON cross_references (source_verse_id, rank)`)

	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE TABLE IF NOT EXISTS dictionary_entries (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		word TEXT NOT NULL,
		word_key TEXT NOT NULL,
		part_of_speech TEXT,
		definitions TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(`CREATE INDEX IF NOT EXISTS dictionary_entries_word_key ON dictionary_entries (word_key)`)

	return err
}
