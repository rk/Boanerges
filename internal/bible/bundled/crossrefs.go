package bundled

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"boanerges/internal/bible/openbible"
	"boanerges/internal/bible/schema"
)

const crossRefMetaKey = "cross_references_v2"

// ImportCrossReferences loads OpenBible cross-reference bundle (idempotent).
func ImportCrossReferences(dbPath, dataPath string) error {
	if _, err := os.Stat(dataPath); err != nil {
		return fmt.Errorf("cross-reference bundle: %w", err)
	}

	mgr, err := schema.Open(dbPath)
	if err != nil {
		return err
	}
	defer mgr.Close()

	if err := mgr.CreateAppSchema(); err != nil {
		return err
	}

	var completed sql.NullString
	if err := mgr.DB().QueryRow(`SELECT completed_at FROM import_meta WHERE key = ?`, crossRefMetaKey).Scan(&completed); err == nil && completed.Valid {
		return nil
	}

	db := mgr.DB()
	if _, err := db.Exec(`DELETE FROM cross_reference_verses`); err != nil {
		return err
	}

	batch := make([]RefRow, 0, 500)
	for id, ref := range openbible.All() {
		batch = append(batch, RefRow{
			ID: id, BookID: ref.BookID, Chapter: ref.Chapter, Verse: ref.Verse,
		})
		if len(batch) >= 500 {
			if err := insertVerseBatch(db, batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		if err := insertVerseBatch(db, batch); err != nil {
			return err
		}
	}

	if _, err := db.Exec(`DELETE FROM cross_references`); err != nil {
		return err
	}

	f, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	refs := make([]crossRefRow, 0, 1000)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "From Verse") {
			continue
		}
		parts := strings.Split(line, "\t")
		if len(parts) < 3 {
			continue
		}
		source := openbible.VerseIDsFromReference(parts[0])
		target := openbible.VerseIDsFromReference(parts[1])
		if source == nil || target == nil {
			continue
		}
		var endID any
		if target.End != nil {
			endID = *target.End
		}
		refs = append(refs, crossRefRow{
			SourceID: source.Start,
			Rank:     atoi(parts[2]),
			Target:   target.Start,
			End:      endID,
		})
		if len(refs) >= 1000 {
			if err := insertCrossRefBatch(db, refs); err != nil {
				return err
			}
			refs = refs[:0]
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(refs) > 0 {
		if err := insertCrossRefBatch(db, refs); err != nil {
			return err
		}
	}

	_, err = db.Exec(`INSERT INTO import_meta (key, completed_at) VALUES (?, ?)
		ON CONFLICT(key) DO UPDATE SET completed_at = excluded.completed_at`,
		crossRefMetaKey, time.Now().UTC().Format(time.RFC3339))

	return err
}

type RefRow struct {
	ID      int
	BookID  string
	Chapter int
	Verse   int
}

type crossRefRow struct {
	SourceID int
	Rank     int
	Target   int
	End      any
}

func insertVerseBatch(db *sql.DB, batch []RefRow) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO cross_reference_verses (id, book_id, chapter, verse) VALUES (?, ?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, row := range batch {
		if _, err := stmt.Exec(row.ID, row.BookID, row.Chapter, row.Verse); err != nil {
			_ = stmt.Close()
			_ = tx.Rollback()
			return err
		}
	}
	_ = stmt.Close()

	return tx.Commit()
}

func insertCrossRefBatch(db *sql.DB, batch []crossRefRow) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT INTO cross_references (source_verse_id, rank, target_start_id, target_end_id) VALUES (?, ?, ?, ?)`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	for _, row := range batch {
		if _, err := stmt.Exec(row.SourceID, row.Rank, row.Target, row.End); err != nil {
			_ = stmt.Close()
			_ = tx.Rollback()
			return err
		}
	}
	_ = stmt.Close()

	return tx.Commit()
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(strings.TrimSpace(s), "%d", &n)

	return n
}
