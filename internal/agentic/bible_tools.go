//go:build agentic

package agentic

import (
	"fmt"
	"strings"

	"boanerges/internal/bible/markup"
	"boanerges/internal/bible/search"
)

func (e *toolEnv) defaultTranslation() string {
	if e.translation != "" {
		return e.translation
	}
	if m, err := loadManifest(e.root); err == nil && strings.TrimSpace(m.Translation) != "" {
		e.translation = m.Translation
		return e.translation
	}
	e.translation = "asv"
	return e.translation
}

func paramInt(params map[string]any, key string) int {
	v, ok := params[key]
	if !ok {
		return 0
	}
	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	default:
		return 0
	}
}

func paramString(params map[string]any, key string) string {
	v, ok := params[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}

type searchArgs struct {
	Query       string `json:"query"`
	Translation string `json:"translation"`
	Limit       int    `json:"limit"`
	PageSize    int    `json:"pageSize"`
	Offset      int    `json:"offset"`
	BookID      string `json:"bookId"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
}

func (a *searchArgs) validate() error {
	if strings.TrimSpace(a.Query) == "" {
		if a.BookID != "" || a.Chapter > 0 || a.Verse > 0 {
			return fmt.Errorf("search needs args.query (keywords, min 2 chars); to load verse text use read_reference with bookId, chapter, verse")
		}
		return fmt.Errorf("search needs args.query (min 2 characters)")
	}
	return nil
}

func (a *searchArgs) pageSize() int {
	if a.PageSize > 0 {
		return a.PageSize
	}
	if a.Limit > 0 {
		return a.Limit
	}
	return 20
}

func (e *toolEnv) keywordSearch(query, translation string, limit, offset int) (search.Response, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	if translation == "" {
		translation = e.defaultTranslation()
	}
	return e.bible.Search(query, translation, limit, offset)
}

type readReferenceArgs struct {
	Translation string `json:"translation"`
	BookID      string `json:"bookId"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
	VerseStart  int    `json:"verseStart"`
	VerseEnd    int    `json:"verseEnd"`
}

func (a *readReferenceArgs) verseRange() (start, end int, err error) {
	start = a.VerseStart
	end = a.VerseEnd
	if start <= 0 && a.Verse > 0 {
		start = a.Verse
	}
	if end <= 0 {
		end = start
	}
	if a.Chapter <= 0 || start <= 0 {
		return 0, 0, fmt.Errorf("read_reference needs bookId, chapter, and verse (or verseStart/verseEnd)")
	}
	if end < start {
		return 0, 0, fmt.Errorf("verseEnd must be >= verseStart")
	}
	return start, end, nil
}

// readReferencePlain returns primary-translation verse text (no HTML), one line per verse.
func (e *toolEnv) readReferencePlain(translation, bookID string, chapter, verseStart, verseEnd int) (string, error) {
	if strings.TrimSpace(bookID) == "" {
		return "", fmt.Errorf("bookId required")
	}
	if translation == "" {
		translation = e.defaultTranslation()
	}
	resp, err := e.bible.Chapter(translation, bookID, chapter)
	if err != nil {
		return "", err
	}
	plain := markup.DefaultFormatter().ToPlainText
	var b strings.Builder
	for _, v := range resp.Chapter.Verses {
		if v.Number < verseStart || v.Number > verseEnd {
			continue
		}
		if b.Len() > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "%d %s", v.Number, strings.TrimSpace(plain(v.Text)))
	}
	if b.Len() == 0 {
		return "", fmt.Errorf("no verses %d–%d in %s %d", verseStart, verseEnd, bookID, chapter)
	}
	return b.String(), nil
}

type agentCrossRef struct {
	BookID   string `json:"bookId"`
	Chapter  int    `json:"chapter"`
	Verse    int    `json:"verse"`
	EndVerse *int   `json:"endVerse,omitempty"`
}

func (e *toolEnv) crossRefsArray(bookID string, chapter, verse int) ([]agentCrossRef, error) {
	if strings.TrimSpace(bookID) == "" || chapter <= 0 || verse <= 0 {
		return nil, fmt.Errorf("cross_references needs bookId, chapter, verse")
	}
	resp, err := e.bible.CrossReferences(bookID, chapter, verse)
	if err != nil {
		return nil, err
	}
	out := make([]agentCrossRef, 0, len(resp.References))
	for _, r := range resp.References {
		out = append(out, agentCrossRef{
			BookID: r.BookID, Chapter: r.Chapter, Verse: r.Verse, EndVerse: r.EndVerse,
		})
	}
	return out, nil
}

func seedReferenceHint(job string, params map[string]any) string {
	if job != "from_reference" {
		return ""
	}
	book := paramString(params, "bookId")
	ch := paramInt(params, "chapter")
	v := paramInt(params, "verse")
	if book == "" || ch <= 0 || v <= 0 {
		return ""
	}
	return fmt.Sprintf(
		"Seed reference: %s %d:%d.\n"+
			"1) read_reference on seed 2) list_topics 3) list every topic in the verse "+
			"4) run_subagent per topic (subagent must upsert_topic) 5) done only after all topics saved.\n"+
			"First tool: {\"tool\":\"read_reference\",\"args\":{\"bookId\":%q,\"chapter\":%d,\"verse\":%d}}\n",
		book, ch, v, book, ch, v,
	)
}
