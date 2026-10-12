package search

import (
	"database/sql"
	"fmt"
	"html"
	"regexp"
	"strings"

	"boanerges/internal/bible"
	"boanerges/internal/bible/schema"
	"boanerges/internal/bible/translations"
)

type Result struct {
	BookID      string `json:"bookId"`
	Chapter     int    `json:"chapter"`
	Verse       int    `json:"verse"`
	Snippet     string `json:"snippet"`
	Translation string `json:"translation"`
}

type Response struct {
	Results []Result `json:"results"`
	Total   int      `json:"total"`
	HasMore bool     `json:"hasMore"`
}

var markStrip = regexp.MustCompile(`</?mark>`)

func Search(db *sql.DB, mgr *schema.Manager, store *translations.Store, query string, translation *string, limit, offset int) (Response, error) {
	if len(strings.TrimSpace(query)) < 2 {
		return Response{}, fmt.Errorf("query too short")
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	var abbrevs []string
	if translation != nil && *translation != "" {
		rec, err := store.FindByAbbrev(*translation)
		if err != nil {
			return Response{}, err
		}
		if rec == nil || !rec.IsReady() {
			return Response{}, fmt.Errorf("translation not found")
		}
		abbrevs = []string{rec.Abbrev}
	} else {
		recs, err := store.ListReady()
		if err != nil {
			return Response{}, err
		}
		for _, r := range recs {
			abbrevs = append(abbrevs, r.Abbrev)
		}
	}

	var results []Result
	for _, abbrev := range abbrevs {
		fts, err := mgr.FtsTable(abbrev)
		if err != nil {
			continue
		}
		if exists, err := mgr.TableExists(fts); err != nil || !exists {
			continue
		}

		rows, err := db.Query(fmt.Sprintf(`
			SELECT book_id, chapter, verse, snippet(%s, 0, '<mark>', '</mark>', '…', 32) AS snippet
			FROM %s WHERE %s MATCH ? LIMIT ? OFFSET ?`, fts, fts, fts), query, limit+1, offset)
		if err != nil {
			continue
		}
		for rows.Next() {
			var bookID int64
			var chapter, verse int
			var snippet string
			if err := rows.Scan(&bookID, &chapter, &verse, &snippet); err != nil {
				rows.Close()
				return Response{}, err
			}
			osis, err := bookOsis(db, mgr, abbrev, bookID)
			if err != nil {
				continue
			}
			canon, _ := bible.NormalizeBookID(osis)
			if canon == "" {
				canon = osis
			}
			results = append(results, Result{
				BookID: canon, Chapter: chapter, Verse: verse,
				Snippet: stripMark(html.UnescapeString(snippet)), Translation: abbrev,
			})
		}
		rows.Close()
	}

	hasMore := len(results) > limit
	page := results
	if hasMore {
		page = results[:limit]
	}
	total := offset + len(page)
	if hasMore {
		total++
	}

	return Response{Results: page, Total: total, HasMore: hasMore}, nil
}

func bookOsis(db *sql.DB, mgr *schema.Manager, abbrev string, bookID int64) (string, error) {
	table, err := mgr.BooksTable(abbrev)
	if err != nil {
		return "", err
	}
	var osis string
	err = db.QueryRow(fmt.Sprintf(`SELECT osis_id FROM %s WHERE id = ?`, table), bookID).Scan(&osis)

	return osis, err
}

func stripMark(s string) string {
	s = markStrip.ReplaceAllString(s, "")

	return strings.TrimSpace(s)
}
