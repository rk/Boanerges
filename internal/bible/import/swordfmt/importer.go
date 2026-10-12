package swordfmt

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"boanerges/internal/bible/catalog"
	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/swordpath"
	"boanerges/sword"
)

type Importer struct {
	schema     *schema.Manager
	formatter  *markup.Formatter
	swordRoots []string
}

func (i *Importer) Format() catalog.ImportFormat {
	return catalog.FormatSword
}

func NewImporter(mgr *schema.Manager, formatter *markup.Formatter, swordRoots []string) *Importer {
	if formatter == nil {
		formatter = markup.DefaultFormatter()
	}

	return &Importer{schema: mgr, formatter: formatter, swordRoots: swordRoots}
}

func (i *Importer) IsModuleAvailable(moduleKey string) bool {
	return swordpath.ModuleInstalled(i.swordRoots, moduleKey)
}

func (i *Importer) ExtractZip(localRoot, zipPath string) error {
	if err := os.MkdirAll(localRoot, 0o755); err != nil {
		return err
	}

	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("open sword zip: %w", err)
	}
	defer r.Close()

	for _, f := range r.File {
		if err := extractFile(localRoot, f); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(dest string, f *zip.File) error {
	name := filepath.Clean(f.Name)
	if strings.HasPrefix(name, "..") {
		return fmt.Errorf("invalid zip path %q", f.Name)
	}
	target := filepath.Join(dest, name)
	if f.FileInfo().IsDir() {
		return os.MkdirAll(target, 0o755)
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	w, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, rc)

	return err
}

func (i *Importer) Import(abbrev, moduleKey string, onProgress func(done, total int)) error {
	confPath, err := swordpath.FindConf(i.swordRoots, moduleKey)
	if err != nil {
		return fmt.Errorf("sword module %q: %w", moduleKey, err)
	}

	mod, err := sword.OpenModule(confPath)
	if err != nil {
		return err
	}

	a, err := i.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}

	booksTable, _ := i.schema.BooksTable(a)
	versesTable, _ := i.schema.VersesTable(a)

	db := i.schema.DB()
	if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", versesTable)); err != nil {
		return err
	}
	if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", booksTable)); err != nil {
		return err
	}

	books := mod.Books()
	total := len(books)
	done := 0

	for _, book := range books {
		res, err := db.Exec(fmt.Sprintf(`
			INSERT INTO %s (name, osis_id, testament, chapters) VALUES (?, ?, ?, ?)`,
			booksTable), book.Name, book.CanonOSIS, book.Testament, book.MaxChapter)
		if err != nil {
			return err
		}
		bookID, _ := res.LastInsertId()

		for ch := 1; ch <= book.MaxChapter; ch++ {
			verseCount := book.Verses[ch-1]
			for v := 1; v <= verseCount; v++ {
				text, err := mod.GetVerse(book.Abbrev, ch, v)
				if err != nil {
					text = ""
				}
				plain := i.formatter.ToPlainText(text)
				if _, err := db.Exec(fmt.Sprintf(`
					INSERT INTO %s (book_id, chapter, verse, text, plain_text) VALUES (?, ?, ?, ?, ?)`,
					versesTable), bookID, ch, v, text, plain); err != nil {
					return err
				}
			}
		}

		done++
		if onProgress != nil {
			onProgress(done, total)
		}
	}

	return nil
}

func (i *Importer) Verify(moduleKey string) error {
	confPath, err := swordpath.FindConf(i.swordRoots, moduleKey)
	if err != nil {
		return err
	}
	mod, err := sword.OpenModule(confPath)
	if err != nil {
		return err
	}

	checks := [][3]any{
		{"Gen", 1, 1},
		{"Matt", 1, 1},
	}
	found := false
	for _, c := range checks {
		text, err := mod.GetVerse(c[0].(string), c[1].(int), c[2].(int))
		if err == nil && strings.TrimSpace(text) != "" {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("verification failed: no reference verses found")
	}

	return nil
}
