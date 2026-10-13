package search

import (
	"database/sql"
	"fmt"
	"html"
	"regexp"
	"sort"
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
		sort.Strings(abbrevs)
	}

	if len(abbrevs) == 1 {
		return searchOne(db, mgr, abbrevs[0], query, limit, offset)
	}

	return searchMany(db, mgr, abbrevs, query, limit, offset)
}

func searchOne(db *sql.DB, mgr *schema.Manager, abbrev, query string, limit, offset int) (Response, error) {
	fts, err := mgr.FtsTable(abbrev)
	if err != nil {
		return Response{}, err
	}
	if exists, err := mgr.TableExists(fts); err != nil || !exists {
		return Response{Results: nil, Total: 0, HasMore: false}, nil
	}

	var total int
	if err := db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s MATCH ?`, fts, fts), query).Scan(&total); err != nil {
		return Response{}, err
	}

	books, err := mgr.BooksTable(abbrev)
	if err != nil {
		return Response{}, err
	}

	rows, err := db.Query(fmt.Sprintf(`
		SELECT b.osis_id, fts.chapter, fts.verse, snippet(%s, 0, '<mark>', '</mark>', '…', 32) AS snippet
		FROM %s fts
		JOIN %s b ON b.id = fts.book_id
		WHERE %s MATCH ?
		ORDER BY rank
		LIMIT ? OFFSET ?`, fts, fts, books, fts), query, limit+1, offset)
	if err != nil {
		return Response{}, err
	}
	defer rows.Close()

	results, err := scanResults(abbrev, rows)
	if err != nil {
		return Response{}, err
	}

	hasMore := len(results) > limit
	page := results
	if hasMore {
		page = results[:limit]
	}

	return Response{Results: page, Total: total, HasMore: hasMore || offset+len(page) < total}, nil
}

func searchMany(db *sql.DB, mgr *schema.Manager, abbrevs []string, query string, limit, offset int) (Response, error) {
	need := offset + limit + 1
	var merged []Result
	total := 0

	for _, abbrev := range abbrevs {
		fts, err := mgr.FtsTable(abbrev)
		if err != nil {
			continue
		}
		if exists, err := mgr.TableExists(fts); err != nil || !exists {
			continue
		}

		var n int
		if err := db.QueryRow(fmt.Sprintf(`SELECT COUNT(*) FROM %s WHERE %s MATCH ?`, fts, fts), query).Scan(&n); err != nil {
			continue
		}
		total += n

		books, err := mgr.BooksTable(abbrev)
		if err != nil {
			continue
		}

		rows, err := db.Query(fmt.Sprintf(`
			SELECT b.osis_id, fts.chapter, fts.verse, snippet(%s, 0, '<mark>', '</mark>', '…', 32) AS snippet
			FROM %s fts
			JOIN %s b ON b.id = fts.book_id
			WHERE %s MATCH ?
			ORDER BY rank
			LIMIT ?`, fts, fts, books, fts), query, need)
		if err != nil {
			continue
		}
		part, err := scanResults(abbrev, rows)
		rows.Close()
		if err != nil {
			return Response{}, err
		}
		merged = append(merged, part...)
	}

	end := offset + limit
	hasMore := len(merged) > end
	if offset > len(merged) {
		return Response{Results: nil, Total: total, HasMore: total > offset}, nil
	}
	page := merged[offset:]
	if len(page) > limit {
		page = page[:limit]
	}

	return Response{Results: page, Total: total, HasMore: hasMore || offset+len(page) < total}, nil
}

func scanResults(abbrev string, rows *sql.Rows) ([]Result, error) {
	var results []Result
	for rows.Next() {
		var osis string
		var chapter, verse int
		var snippet string
		if err := rows.Scan(&osis, &chapter, &verse, &snippet); err != nil {
			return nil, err
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

	return results, rows.Err()
}

func stripMark(s string) string {
	s = markStrip.ReplaceAllString(s, "")

	return strings.TrimSpace(s)
}
