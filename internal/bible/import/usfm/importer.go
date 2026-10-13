package usfm

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"boanerges/internal/bible"
	"boanerges/internal/bible/catalog"
	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
)

type Importer struct {
	schema    *schema.Manager
	formatter *markup.Formatter
}

func NewImporter(mgr *schema.Manager, formatter *markup.Formatter) *Importer {
	if formatter == nil {
		formatter = markup.DefaultFormatter()
	}

	return &Importer{schema: mgr, formatter: formatter}
}

type importState struct {
	bookDBID      sql.NullInt64
	pendingOsisID string
	chapter       int
	verse         int
	verseText     string
	maxChapter    int
	chapterRows   []schema.VerseRow
}

func (i *Importer) ImportFromZip(abbrev, zipPath string) error {
	dir, err := os.MkdirTemp("", "boanerges-usfm-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open usfm zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		name := filepath.Base(f.Name)
		if strings.HasPrefix(name, "._") || !strings.HasSuffix(strings.ToLower(name), ".usfm") {
			continue
		}
		if err := extractZipFile(f, dir); err != nil {
			return err
		}
	}

	files, err := findUsfmFiles(dir)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no USFM file found in archive")
	}

	a, err := i.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}

	booksTable, _ := i.schema.BooksTable(a)
	versesTable, _ := i.schema.VersesTable(a)

	if _, err := i.schema.DB().Exec(fmt.Sprintf("DELETE FROM %s", versesTable)); err != nil {
		return err
	}
	if _, err := i.schema.DB().Exec(fmt.Sprintf("DELETE FROM %s", booksTable)); err != nil {
		return err
	}

	for _, file := range files {
		if err := i.importFromFile(a, file); err != nil {
			return err
		}
	}

	return nil
}

func (i *Importer) ImportFromFile(abbrev, path string, clearTables bool) error {
	a, err := i.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}

	booksTable, _ := i.schema.BooksTable(a)
	versesTable, _ := i.schema.VersesTable(a)

	if clearTables {
		if _, err := i.schema.DB().Exec(fmt.Sprintf("DELETE FROM %s", versesTable)); err != nil {
			return err
		}
		if _, err := i.schema.DB().Exec(fmt.Sprintf("DELETE FROM %s", booksTable)); err != nil {
			return err
		}
	}

	return i.importFromFile(a, path)
}

func (i *Importer) importFromFile(abbrev, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read usfm: %w", err)
	}

	booksTable, _ := i.schema.BooksTable(abbrev)

	return i.importContent(abbrev, string(content), booksTable)
}

func (i *Importer) Format() catalog.ImportFormat {
	return catalog.FormatUsfm
}

var (
	idMarker      = regexp.MustCompile(`^\\id\s+(\S+)`)
	headingMarker = regexp.MustCompile(`^\\h\s+(.+)$`)
	chapterMarker = regexp.MustCompile(`^\\c\s+(\d+)`)
	verseMarker   = regexp.MustCompile(`^\\v\s+(\d+)(?:-\d+)?\s*(.*)$`)
	markerLine    = regexp.MustCompile(`^\\[a-z0-9*]+\s*(.*)$`)
	usfmPlainRe   = regexp.MustCompile(`\\[a-z0-9*]+(\s+|$)`)
)

func (i *Importer) importContent(abbrev, content, booksTable string) error {
	state := importState{}
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}

		if m := idMarker.FindStringSubmatch(line); len(m) == 2 {
			if err := i.flushBook(&state, abbrev, booksTable); err != nil {
				return err
			}
			osis, ok := bible.NormalizeBookID(m[1])
			if !ok {
				state.pendingOsisID = ""
				continue
			}
			state.pendingOsisID = osis
			continue
		}

		if state.pendingOsisID == "" {
			continue
		}

		if !state.bookDBID.Valid && headingMarker.MatchString(line) {
			m := headingMarker.FindStringSubmatch(line)
			id, err := i.insertBook(booksTable, state.pendingOsisID, strings.TrimSpace(m[1]))
			if err != nil {
				return err
			}
			state.bookDBID = sql.NullInt64{Int64: id, Valid: true}
			continue
		}

		if m := chapterMarker.FindStringSubmatch(line); len(m) == 2 {
			if err := i.ensureBook(&state, booksTable); err != nil {
				return err
			}
			if err := i.flushVerse(&state); err != nil {
				return err
			}
			if err := i.flushChapter(&state, abbrev); err != nil {
				return err
			}
			state.chapter = atoi(m[1])
			state.verse = 0
			state.verseText = ""
			if state.chapter > state.maxChapter {
				state.maxChapter = state.chapter
			}
			continue
		}

		if !state.bookDBIdValid() {
			continue
		}

		if m := verseMarker.FindStringSubmatch(line); len(m) == 3 {
			if err := i.flushVerse(&state); err != nil {
				return err
			}
			state.verse = atoi(m[1])
			state.verseText = m[2]
			continue
		}

		if state.verse < 1 || state.chapter < 1 {
			continue
		}

		if strings.HasPrefix(line, `\`) {
			if m := markerLine.FindStringSubmatch(line); len(m) == 2 && strings.TrimSpace(m[1]) != "" {
				state.verseText = strings.TrimSpace(state.verseText + " " + m[1])
			}
			continue
		}

		state.verseText = strings.TrimSpace(state.verseText + " " + line)
	}

	return i.flushBook(&state, abbrev, booksTable)
}

func (s *importState) bookDBIdValid() bool {
	return s.bookDBID.Valid
}

func (i *Importer) flushVerse(state *importState) error {
	if !state.bookDBID.Valid || state.chapter < 1 || state.verse < 1 || strings.TrimSpace(state.verseText) == "" {
		return nil
	}

	text := strings.TrimSpace(state.verseText)
	plain := usfmPlainText(text)
	if plain == "" {
		plain = i.formatter.ToPlainText(text)
	}

	state.chapterRows = append(state.chapterRows, schema.VerseRow{
		BookID: state.bookDBID.Int64, Chapter: state.chapter, Verse: state.verse,
		Text: text, PlainText: plain,
	})

	return nil
}

func (i *Importer) flushChapter(state *importState, abbrev string) error {
	if len(state.chapterRows) == 0 {
		return nil
	}
	if err := i.schema.InsertVerses(abbrev, state.chapterRows); err != nil {
		return err
	}
	state.chapterRows = state.chapterRows[:0]

	return nil
}

func (i *Importer) flushBook(state *importState, abbrev, booksTable string) error {
	if err := i.flushVerse(state); err != nil {
		return err
	}
	if err := i.flushChapter(state, abbrev); err != nil {
		return err
	}
	if state.bookDBID.Valid {
		_, _ = i.schema.DB().Exec(fmt.Sprintf("UPDATE %s SET chapters = ? WHERE id = ?", booksTable),
			state.maxChapter, state.bookDBID.Int64)
	}
	state.bookDBID = sql.NullInt64{}
	state.pendingOsisID = ""
	state.chapter = 0
	state.verse = 0
	state.verseText = ""
	state.maxChapter = 0
	state.chapterRows = nil

	return nil
}

func (i *Importer) ensureBook(state *importState, booksTable string) error {
	if state.bookDBID.Valid || state.pendingOsisID == "" {
		return nil
	}
	id, err := i.insertBook(booksTable, state.pendingOsisID, strings.ToUpper(state.pendingOsisID))
	if err != nil {
		return err
	}
	state.bookDBID = sql.NullInt64{Int64: id, Valid: true}

	return nil
}

func (i *Importer) insertBook(booksTable, osisID, name string) (int64, error) {
	res, err := i.schema.DB().Exec(fmt.Sprintf(`
		INSERT INTO %s (name, osis_id, testament, chapters) VALUES (?, ?, ?, 0)`, booksTable),
		name, osisID, guessTestament(osisID))
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func guessTestament(osisID string) string {
	nt := map[string]struct{}{
		"mat": {}, "mrk": {}, "luk": {}, "jhn": {}, "act": {}, "rom": {}, "1co": {}, "2co": {},
		"gal": {}, "eph": {}, "php": {}, "col": {}, "1th": {}, "2th": {}, "1ti": {}, "2ti": {},
		"tit": {}, "phm": {}, "heb": {}, "jas": {}, "1pe": {}, "2pe": {}, "1jn": {}, "2jn": {},
		"3jn": {}, "jud": {}, "rev": {},
	}
	if _, ok := nt[strings.ToLower(osisID)]; ok {
		return "nt"
	}

	return "ot"
}

func usfmPlainText(text string) string {
	plain := usfmPlainRe.ReplaceAllString(text, " ")
	plain = strings.Join(strings.Fields(plain), " ")

	return strings.TrimSpace(plain)
}

func findUsfmFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		name := d.Name()
		if strings.HasPrefix(name, "._") || !strings.HasSuffix(strings.ToLower(name), ".usfm") {
			return nil
		}
		files = append(files, path)

		return nil
	})
	sort.Slice(files, func(a, b int) bool { return files[a] < files[b] })

	return files, err
}

func extractZipFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	outPath := filepath.Join(dest, filepath.Base(f.Name))
	w, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, rc)

	return err
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)

	return n
}
