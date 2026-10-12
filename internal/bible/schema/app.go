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

	return err
}
