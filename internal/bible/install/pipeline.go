package install

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"boanerges/internal/bible/progress"
	"boanerges/internal/bible/catalog"
	bibleimport "boanerges/internal/bible/import"
	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/swordpath"
	"boanerges/internal/bible/translations"
	"boanerges/sword"
)

type Pipeline struct {
	catalog     *catalog.Catalog
	store       *translations.Store
	schema      *schema.Manager
	importers   *bibleimport.Registry
	formatter   *markup.Formatter
	swordRoots  []string
	localSword  string
	tmpDir      string
	progress    *progress.Emitter
	httpClient  *http.Client
}

func NewPipeline(
	cat *catalog.Catalog,
	store *translations.Store,
	mgr *schema.Manager,
	importers *bibleimport.Registry,
	bundledSword, localSword, tmpDir string,
	emitter *progress.Emitter,
) *Pipeline {
	roots := swordpath.Roots(bundledSword, localSword)

	return &Pipeline{
		catalog:    cat,
		store:      store,
		schema:     mgr,
		importers:  importers,
		formatter:  markup.DefaultFormatter(),
		swordRoots: roots,
		localSword: localSword,
		tmpDir:     tmpDir,
		progress:   emitter,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
}

func (p *Pipeline) Run(translationID int64) error {
	rec, err := p.store.FindByID(translationID)
	if err != nil || rec == nil {
		return fmt.Errorf("translation not found")
	}
	if rec.IsReady() {
		return nil
	}

	entry, err := p.catalog.Find(rec.Abbrev)
	if err != nil {
		entry, err = p.catalog.FindByID(rec.Abbrev)
	}
	if err != nil {
		return err
	}

	var downloaded string
	defer func() {
		if downloaded != "" {
			_ = os.Remove(downloaded)
		}
	}()

	try := func(step string, fn func() error) error {
		start := time.Now()
		err := fn()
		log.Printf("translation install %s: %s in %s", rec.Abbrev, step, time.Since(start).Round(time.Millisecond))
		if err != nil {
			_ = p.store.MarkFailed(translationID, err.Error())
			p.emit(rec.Abbrev, translations.StatusFailed, "failed", err.Error())

			return err
		}

		return nil
	}

	resume := p.importLooksComplete(rec.Abbrev)
	if resume {
		log.Printf("translation install %s: resume (verses already in DB)", rec.Abbrev)
	} else {
		if err := try("download", func() error { return p.download(rec, entry, &downloaded) }); err != nil {
			return err
		}
		if err := try("create_schema", func() error { return p.createSchema(rec.Abbrev) }); err != nil {
			return err
		}
		if err := try("import_verses", func() error { return p.importVerses(rec, entry, downloaded) }); err != nil {
			return err
		}
		if err := try("verify", func() error { return p.verify(rec, entry) }); err != nil {
			return err
		}
	}

	if err := try("index_coverage", func() error { return p.schema.IndexBookCoverage(rec.Abbrev) }); err != nil {
		return err
	}
	if err := try("backfill_plain_text", func() error {
		return p.schema.BackfillPlainTextIfEmpty(rec.Abbrev, p.formatter)
	}); err != nil {
		return err
	}

	ftsOk, err := p.schema.FtsPopulated(rec.Abbrev)
	if err != nil {
		return err
	}
	if !ftsOk {
		if err := try("build_fts", func() error { return p.buildFTS(rec.Abbrev) }); err != nil {
			return err
		}
	}

	if err := p.store.MarkReady(translationID); err != nil {
		return err
	}
	p.emit(rec.Abbrev, translations.StatusReady, "ready", "")

	return nil
}

func (p *Pipeline) importLooksComplete(abbrev string) bool {
	n, err := p.schema.VerseCount(abbrev)
	if err != nil || n < 1000 {
		return false
	}

	return p.hasVerse(abbrev, "gen", 1, 1) || p.hasVerse(abbrev, "mat", 1, 1)
}

func (p *Pipeline) download(rec *translations.Record, entry catalog.Entry, downloaded *string) error {
	switch entry.ImportAs {
	case catalog.FormatSword:
		if p.importers.Sword().IsModuleAvailable(entry.Short) {
			p.updateProgress(rec.ID, rec.Abbrev, translations.StatusDownloading, "source_ready")
			return p.syncMetadata(rec.ID, entry)
		}
	case catalog.FormatUsfm:
		// always download USFM
	}

	p.updateProgress(rec.ID, rec.Abbrev, translations.StatusDownloading, "downloading")

	if err := os.MkdirAll(p.tmpDir, 0o755); err != nil {
		return err
	}

	zipPath := filepath.Join(p.tmpDir, entry.Short+".module.zip")
	if err := p.downloadURL(entry.URL, zipPath); err != nil {
		return err
	}

	switch entry.ImportAs {
	case catalog.FormatSword:
		if err := p.importers.Sword().ExtractZip(p.localSword, zipPath); err != nil {
			return err
		}
		_ = os.Remove(zipPath)
	case catalog.FormatUsfm:
		*downloaded = zipPath
	default:
		return fmt.Errorf("unsupported import format")
	}

	p.updateProgress(rec.ID, rec.Abbrev, translations.StatusDownloading, "downloaded")

	return p.syncMetadata(rec.ID, entry)
}

func (p *Pipeline) downloadURL(url, dest string) error {
	resp, err := p.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to download: HTTP %d", resp.StatusCode)
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, resp.Body)
	if err != nil {
		_ = os.Remove(dest)
		return err
	}
	if n == 0 {
		_ = os.Remove(dest)
		return fmt.Errorf("failed to download: empty file")
	}

	return nil
}

func (p *Pipeline) syncMetadata(id int64, entry catalog.Entry) error {
	format := entry.MarkupFormat
	if format == "" && entry.ImportAs == catalog.FormatUsfm {
		format = "usfm"
	}

	if entry.ImportAs == catalog.FormatSword {
		if conf, err := swordpath.FindConf(p.swordRoots, entry.Short); err == nil {
			_, meta, err := sword.ParseConfFile(conf)
			if err == nil {
				return p.store.ApplySwordMetadata(id, meta.Name, meta.Format, meta.Versification, meta.About,
					meta.VersionString, meta.VersionDate, meta.Copyright, meta.CopyrightContact, meta.Source)
			}
		}
	}

	return p.store.UpdateCatalogMeta(id, entry.Name, format, entry.About)
}

func (p *Pipeline) createSchema(abbrev string) error {
	a, err := p.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}
	if err := p.schema.DropTables(a); err != nil {
		return err
	}
	if err := p.schema.CreateTables(a); err != nil {
		return err
	}

	rec, _ := p.store.FindByAbbrev(a)
	if rec != nil {
		p.updateProgress(rec.ID, a, translations.StatusCreatingSchema, "creating_schema")
	}

	return nil
}

func (p *Pipeline) importVerses(rec *translations.Record, entry catalog.Entry, downloaded string) error {
	p.updateProgress(rec.ID, rec.Abbrev, translations.StatusImporting, "importing")

	switch entry.ImportAs {
	case catalog.FormatSword:
		return p.importers.Sword().Import(rec.Abbrev, entry.Short, func(done, total int) {
			if total > 0 {
				p.emit(rec.Abbrev, translations.StatusImporting, "importing", "")
			}
		})
	case catalog.FormatUsfm:
		if downloaded == "" {
			return fmt.Errorf("missing USFM source")
		}
		return p.importers.Usfm().ImportFromZip(rec.Abbrev, downloaded)
	default:
		return fmt.Errorf("unsupported import format")
	}
}

func (p *Pipeline) verify(rec *translations.Record, entry catalog.Entry) error {
	switch entry.ImportAs {
	case catalog.FormatSword:
		if err := p.importers.Sword().Verify(entry.Short); err != nil {
			return err
		}
	case catalog.FormatUsfm:
		if !p.hasVerse(rec.Abbrev, "gen", 1, 1) && !p.hasVerse(rec.Abbrev, "mat", 1, 1) {
			return fmt.Errorf("verification failed: no reference verses found")
		}
	}

	p.updateProgress(rec.ID, rec.Abbrev, translations.StatusVerifying, "verifying")

	return nil
}

func (p *Pipeline) hasVerse(abbrev, book string, chapter, verse int) bool {
	a, err := p.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return false
	}
	booksTable, _ := p.schema.BooksTable(a)
	versesTable, _ := p.schema.VersesTable(a)

	var bookID int64
	err = p.schema.DB().QueryRow(fmt.Sprintf("SELECT id FROM %s WHERE osis_id = ?", booksTable), strings.ToLower(book)).Scan(&bookID)
	if err != nil {
		return false
	}

	var text string
	err = p.schema.DB().QueryRow(fmt.Sprintf(`
		SELECT text FROM %s WHERE book_id = ? AND chapter = ? AND verse = ? AND text != ''`,
		versesTable), bookID, chapter, verse).Scan(&text)

	return err == nil
}

func (p *Pipeline) buildFTS(abbrev string) error {
	a, err := p.schema.ValidateAbbrev(abbrev)
	if err != nil {
		return err
	}

	rec, _ := p.store.FindByAbbrev(a)
	if rec != nil {
		p.updateProgress(rec.ID, a, translations.StatusIndexing, "indexing")
	}

	if err := p.schema.RebuildFtsIndex(a, p.formatter); err != nil {
		return err
	}

	if rec != nil {
		p.updateProgress(rec.ID, a, translations.StatusIndexing, "indexed")
	}

	return nil
}

func (p *Pipeline) updateProgress(id int64, abbrev, status, step string) {
	_ = p.store.UpdateProgress(id, status, step, nil)
	p.emit(abbrev, status, step, "")
}

func (p *Pipeline) emit(abbrev, status, step, errMsg string) {
	if p.progress == nil {
		return
	}

	p.progress.Emit(progress.Payload{
		Abbrev:        abbrev,
		Step:          step,
		Percent:       translations.ProgressPercent(status, step),
		Error:         errMsg,
		InstallStatus: status,
	})
}
