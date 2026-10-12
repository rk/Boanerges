package translations

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"boanerges/internal/bible/schema"
)

type Store struct {
	mgr *schema.Manager
}

func NewStore(mgr *schema.Manager) *Store {
	return &Store{mgr: mgr}
}

func (s *Store) DB() *sql.DB {
	return s.mgr.DB()
}

func (s *Store) EnsureSchema() error {
	return s.mgr.CreateAppSchema()
}

func (s *Store) FindByAbbrev(abbrev string) (*Record, error) {
	a, err := s.mgr.ValidateAbbrev(abbrev)
	if err != nil {
		a = strings.ToLower(abbrev)
	}

	row := s.mgr.DB().QueryRow(`
		SELECT id, abbrev, name, COALESCE(format,''), COALESCE(about,''),
			install_status, COALESCE(install_step,''), COALESCE(install_error,''), bundled
		FROM translations WHERE abbrev = ?`, a)

	var rec Record
	var bundled int
	err = row.Scan(
		&rec.ID, &rec.Abbrev, &rec.Name, &rec.Format, &rec.About,
		&rec.InstallStatus, &rec.InstallStep, &rec.InstallError, &bundled,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rec.Bundled = bundled != 0

	return &rec, nil
}

func (s *Store) FindByID(id int64) (*Record, error) {
	row := s.mgr.DB().QueryRow(`
		SELECT id, abbrev, name, COALESCE(format,''), COALESCE(about,''),
			install_status, COALESCE(install_step,''), COALESCE(install_error,''), bundled
		FROM translations WHERE id = ?`, id)

	var rec Record
	var bundled int
	err := row.Scan(
		&rec.ID, &rec.Abbrev, &rec.Name, &rec.Format, &rec.About,
		&rec.InstallStatus, &rec.InstallStep, &rec.InstallError, &bundled,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	rec.Bundled = bundled != 0

	return &rec, nil
}

func (s *Store) ListReady() ([]Record, error) {
	rows, err := s.mgr.DB().Query(`
		SELECT id, abbrev, name, COALESCE(format,''), COALESCE(about,''),
			install_status, COALESCE(install_step,''), COALESCE(install_error,''), bundled
		FROM translations WHERE install_status = ? ORDER BY name`, StatusReady)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Record
	for rows.Next() {
		var rec Record
		var bundled int
		if err := rows.Scan(
			&rec.ID, &rec.Abbrev, &rec.Name, &rec.Format, &rec.About,
			&rec.InstallStatus, &rec.InstallStep, &rec.InstallError, &bundled,
		); err != nil {
			return nil, err
		}
		rec.Bundled = bundled != 0
		out = append(out, rec)
	}

	return out, rows.Err()
}

type CreateParams struct {
	Abbrev   string
	Name     string
	Format   string
	About    string
	Bundled  bool
}

func (s *Store) Create(p CreateParams) (*Record, error) {
	a, err := s.mgr.ValidateAbbrev(p.Abbrev)
	if err != nil {
		return nil, err
	}

	bundled := 0
	if p.Bundled {
		bundled = 1
	}

	res, err := s.mgr.DB().Exec(`
		INSERT INTO translations (abbrev, name, format, about, install_status, install_step, bundled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a, p.Name, nullIfEmpty(p.Format), nullIfEmpty(p.About),
		StatusPending, "pending", bundled, now(), now(),
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()

	return s.FindByID(id)
}

func (s *Store) UpdateCatalogMeta(id int64, name, format, about string) error {
	_, err := s.mgr.DB().Exec(`
		UPDATE translations SET name = ?, format = COALESCE(NULLIF(?, ''), format),
			about = COALESCE(NULLIF(?, ''), about), updated_at = ? WHERE id = ?`,
		name, format, about, now(), id,
	)

	return err
}

func (s *Store) UpdateProgress(id int64, status, step string, errMsg *string) error {
	var installErr any
	if errMsg != nil {
		installErr = *errMsg
	}
	_, err := s.mgr.DB().Exec(`
		UPDATE translations SET install_status = ?, install_step = ?, install_error = ?, updated_at = ? WHERE id = ?`,
		status, step, installErr, now(), id,
	)

	return err
}

func (s *Store) MarkFailed(id int64, msg string) error {
	return s.UpdateProgress(id, StatusFailed, "failed", &msg)
}

func (s *Store) MarkReady(id int64) error {
	return s.UpdateProgress(id, StatusReady, "ready", nil)
}

func (s *Store) Delete(id int64) error {
	_, err := s.mgr.DB().Exec(`DELETE FROM translations WHERE id = ?`, id)

	return err
}

func (s *Store) ApplySwordMetadata(id int64, name, format, versification, about, versionString, versionDate, copyright, contact, source string) error {
	_, err := s.mgr.DB().Exec(`
		UPDATE translations SET
			name = COALESCE(NULLIF(?, ''), name),
			format = COALESCE(NULLIF(?, ''), format),
			versification = COALESCE(NULLIF(?, ''), versification),
			about = COALESCE(NULLIF(?, ''), about),
			version_string = COALESCE(NULLIF(?, ''), version_string),
			version_date = COALESCE(NULLIF(?, ''), version_date),
			copyright = COALESCE(NULLIF(?, ''), copyright),
			copyright_contact = COALESCE(NULLIF(?, ''), copyright_contact),
			source = COALESCE(NULLIF(?, ''), source),
			updated_at = ?
		WHERE id = ?`,
		name, format, versification, about, versionString, versionDate, copyright, contact, source, now(), id,
	)

	return err
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}

	return s
}

func now() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05")
}

// ProgressPercent mirrors legacy install progress for API responses.
func ProgressPercent(status, step string) int {
	if status == StatusReady {
		return 100
	}
	switch step {
	case "ready":
		return 100
	case "indexing", "indexed":
		return 85
	case "importing":
		return 50
	case "creating_schema":
		return 30
	case "downloaded", "source_ready":
		return 10
	default:
		return 0
	}
}

func (s *Store) MustAbbrev(id int64) (string, error) {
	rec, err := s.FindByID(id)
	if err != nil {
		return "", err
	}
	if rec == nil {
		return "", fmt.Errorf("translation id %d not found", id)
	}

	return rec.Abbrev, nil
}
