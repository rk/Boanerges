package studyprint

import (
	"fmt"
	"strings"

	"boanerges/internal/bible"
	"boanerges/internal/bible/read"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/translations"
	"boanerges/internal/settings"
	"boanerges/internal/studydata"
)

type printColumn struct {
	Label   string
	Kind    string
	Verses  []read.Verse
	Content string
	Message string
	Entries []verseListEntry
}

type verseListEntry struct {
	Label string
	Text  string
}

func buildColumns(
	reader *read.Reader,
	trans *translations.Store,
	studySvc *studydata.Service,
	study settings.StudySettings,
	primary read.Chapter,
	includeUserWork bool,
) ([]printColumn, error) {
	abbrev := func(id string) string {
		rec, err := trans.FindByAbbrev(id)
		if err != nil || rec == nil {
			return strings.ToUpper(id)
		}

		return strings.ToUpper(rec.Abbrev)
	}

	cols := []printColumn{{
		Label:  fmt.Sprintf("%s %d (%s)", primary.Book, primary.Chapter, abbrev(study.TranslationID)),
		Kind:   "bible",
		Verses: primary.Verses,
	}}

	for slot, colType := range study.Columns {
		switch colType {
		case settings.ColumnBibleSecondary:
			tid := study.TranslationBID
			if slot == 1 {
				tid = study.TranslationCID
			}
			ch, err := reader.Chapter(tid, study.BookID, study.Chapter)
			if err != nil {
				return nil, err
			}
			cols = append(cols, printColumn{
				Label:  fmt.Sprintf("%s %d (%s)", ch.Book, ch.Chapter, abbrev(tid)),
				Kind:   "bible",
				Verses: ch.Verses,
			})
		case settings.ColumnNotes:
			label := fmt.Sprintf("%s %d (Notes)", primary.Book, primary.Chapter)
			if includeUserWork && studySvc != nil {
				text, _ := studySvc.NotesGet(study.BookID, study.Chapter)
				cols = append(cols, printColumn{Label: label, Kind: "notes", Content: text})
			} else {
				cols = append(cols, printColumn{Label: label, Kind: "lined-notes"})
			}
		case settings.ColumnScribe:
			label := fmt.Sprintf("%s %d (Scribe)", primary.Book, primary.Chapter)
			if !includeUserWork || studySvc == nil {
				cols = append(cols, printColumn{Label: label, Kind: "scribe"})
				break
			}
			draft, _ := studySvc.ScribeGet(study.BookID, study.Chapter)
			cols = append(cols, printColumn{
				Label:  label,
				Kind:   "scribe-content",
				Verses: scribeVerses(primary.Verses, draft),
			})
		case settings.ColumnVerseList:
			cols = append(cols, verseListColumn(reader, study, primary))
		case settings.ColumnSearch, settings.ColumnCrossReferences, settings.ColumnComparison, settings.ColumnDictionary:
			cols = append(cols, printColumn{
				Label:   columnLabel(colType),
				Kind:    "message",
				Message: "Interactive view — not included in print.",
			})
		}
	}

	return cols, nil
}

func columnLabel(colType string) string {
	switch colType {
	case settings.ColumnBibleSecondary:
		return "Translation"
	case settings.ColumnNotes:
		return "Notes"
	case settings.ColumnScribe:
		return "Scribe"
	case settings.ColumnSearch:
		return "Search"
	case settings.ColumnCrossReferences:
		return "Cross References"
	case settings.ColumnDictionary:
		return "Dictionary"
	case settings.ColumnComparison:
		return "Comparison"
	case settings.ColumnVerseList:
		return "Verse List"
	default:
		return colType
	}
}

func scribeVerses(source []read.Verse, draft []studydata.ScribeVerse) []read.Verse {
	byVerse := map[int]studydata.ScribeVerse{}
	for _, d := range draft {
		byVerse[d.Verse] = d
	}

	var out []read.Verse
	pendingBreak := false
	for _, src := range source {
		entry, ok := byVerse[src.Number]
		text := ""
		if ok {
			text = strings.TrimSpace(entry.Text)
		}
		startsPara := pendingBreak || (ok && entry.ParagraphStart) || (!ok && src.Number == 1)
		pendingBreak = false
		if text == "" {
			if startsPara && len(out) > 0 {
				pendingBreak = true
			}
			continue
		}
		para := len(out) == 0 || startsPara
		out = append(out, read.Verse{Number: src.Number, Text: text, ParagraphStart: para})
	}

	return out
}

func verseListColumn(reader *read.Reader, study settings.StudySettings, primary read.Chapter) printColumn {
	title := "Verse List"
	var entries []verseListEntry
	if study.VerseListActiveID != nil && *study.VerseListActiveID != "" {
		svc, err := studydata.NewService()
		if err == nil {
			list, err := svc.VerseListGet(*study.VerseListActiveID)
			if err == nil {
				title = list.Title
				for _, e := range list.Entries {
					text := "—"
					if ch, err := reader.Chapter(study.TranslationID, e.BookID, e.Chapter); err == nil {
						for _, v := range ch.Verses {
							if v.Number == e.Verse {
								text = v.Text
								break
							}
						}
					}
					entries = append(entries, verseListEntry{
						Label: fmt.Sprintf("%s %d:%d", bible.DisplayName(e.BookID), e.Chapter, e.Verse),
						Text:  text,
					})
				}
			}
		}
	}

	return printColumn{Label: title, Kind: "verse-list", Entries: entries}
}

func openStudyServices(dbPath string) (*read.Reader, *translations.Store, *schema.Manager, error) {
	mgr, err := schema.Open(dbPath)
	if err != nil {
		return nil, nil, nil, err
	}

	return read.NewReader(mgr.DB(), mgr, nil), translations.NewStore(mgr), mgr, nil
}
