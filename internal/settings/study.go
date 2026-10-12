package settings

import (
	"boanerges/internal/bible"
)

const studyKey = "study"

// StudySettings match legacy study store shape.
type StudySettings struct {
	ColumnCount          int      `json:"columnCount"`
	Columns              []string `json:"columns"`
	BookID               string   `json:"bookId"`
	Chapter              int      `json:"chapter"`
	TranslationID        string   `json:"translationId"`
	TranslationBID       string   `json:"translationBId"`
	TranslationCID       string   `json:"translationCId"`
	VerseListShowContent bool     `json:"verseListShowContent"`
	VerseListActiveID    *string  `json:"verseListActiveId"`
}

type StudyStore struct {
	repo *Repository
}

func NewStudyStore(repo *Repository) *StudyStore {
	return &StudyStore{repo: repo}
}

func (s *StudyStore) Defaults() StudySettings {
	return StudySettings{
		ColumnCount:          1,
		Columns:              []string{},
		BookID:               "gen",
		Chapter:              15,
		TranslationID:        "asv",
		TranslationBID:       "asv",
		TranslationCID:       "asv",
		VerseListShowContent: true,
		VerseListActiveID:    nil,
	}
}

func (s *StudyStore) Get() StudySettings {
	raw := s.repo.Get(studyKey, map[string]any{})
	out := s.migrateFromLegacy(raw)
	if canonical, ok := bible.NormalizeBookID(out.BookID); ok {
		out.BookID = canonical
	}

	return out
}

func (s *StudyStore) Update(next StudySettings) StudySettings {
	if canonical, ok := bible.NormalizeBookID(next.BookID); ok {
		next.BookID = canonical
	}

	merged := s.migrateFromLegacy(studyToMap(next))
	stored := s.repo.Update(studyKey, studyToMap(merged), studyDefaultsMap())

	return mapToStudy(stored)
}

func (s *StudyStore) migrateFromLegacy(stored map[string]any) StudySettings {
	defaults := s.Defaults()

	if _, hasColumnCount := stored["columnCount"]; hasColumnCount {
		columnCount := asInt(stored, "columnCount", 1)
		if columnCount != 1 && columnCount != 2 && columnCount != 3 {
			columnCount = 1
		}

		columns := asStringSlice(stored, "columns")
		if columns == nil {
			columns = s.normalizeColumnSlots(columnCount, []string{})
		} else {
			columns = s.sanitizeColumns(columnCount, columns)
		}

		return mapToStudy(map[string]any{
			"columnCount":          columnCount,
			"columns":              stringSliceToAny(columns),
			"bookId":               asString(stored, "bookId", defaults.BookID),
			"chapter":              asInt(stored, "chapter", defaults.Chapter),
			"translationId":        asString(stored, "translationId", defaults.TranslationID),
			"translationBId":       asString(stored, "translationBId", defaults.TranslationBID),
			"translationCId":       asString(stored, "translationCId", defaults.TranslationCID),
			"verseListShowContent": asBool(stored, "verseListShowContent", defaults.VerseListShowContent),
			"verseListActiveId":    stored["verseListActiveId"],
		})
	}

	activeView := asString(stored, "activeView", "bible")
	layout := legacyLayout(activeView)

	return mapToStudy(map[string]any{
		"columnCount":          layout.columnCount,
		"columns":              stringSliceToAny(layout.columns),
		"bookId":               asString(stored, "bookId", defaults.BookID),
		"chapter":              asInt(stored, "chapter", defaults.Chapter),
		"translationId":        asString(stored, "translationId", defaults.TranslationID),
		"translationBId":       asString(stored, "translationBId", defaults.TranslationBID),
		"translationCId":       asString(stored, "translationCId", defaults.TranslationCID),
		"verseListShowContent": asBool(stored, "verseListShowContent", defaults.VerseListShowContent),
		"verseListActiveId":    stored["verseListActiveId"],
	})
}

type legacyColumnLayout struct {
	columnCount int
	columns     []string
}

func legacyLayout(activeView string) legacyColumnLayout {
	switch activeView {
	case "comparison":
		return legacyColumnLayout{columnCount: 2, columns: []string{ColumnBibleSecondary}}
	case "scribe":
		return legacyColumnLayout{columnCount: 3, columns: []string{ColumnScribe, ColumnBibleSecondary}}
	default:
		return legacyColumnLayout{columnCount: 1, columns: []string{}}
	}
}

func (s *StudyStore) sanitizeColumns(columnCount int, columns []string) []string {
	filtered := make([]string, 0, len(columns))
	for _, column := range columns {
		if allowedColumnTypes[column] {
			filtered = append(filtered, column)
		}
	}

	return s.normalizeColumnSlots(columnCount, filtered)
}

func (s *StudyStore) normalizeColumnSlots(columnCount int, columns []string) []string {
	needed := max(0, columnCount-1)
	normalized := columns
	if len(normalized) > needed {
		normalized = normalized[:needed]
	}
	for len(normalized) < needed {
		normalized = append(normalized, ColumnBibleSecondary)
	}

	return normalized
}

func studyDefaultsMap() map[string]any {
	d := StudySettings{
		ColumnCount:          1,
		Columns:              []string{},
		BookID:               "gen",
		Chapter:              15,
		TranslationID:        "asv",
		TranslationBID:       "asv",
		TranslationCID:       "asv",
		VerseListShowContent: true,
		VerseListActiveID:    nil,
	}

	return studyToMap(d)
}

func studyToMap(s StudySettings) map[string]any {
	m := map[string]any{
		"columnCount":          s.ColumnCount,
		"columns":              stringSliceToAny(s.Columns),
		"bookId":               s.BookID,
		"chapter":              s.Chapter,
		"translationId":        s.TranslationID,
		"translationBId":       s.TranslationBID,
		"translationCId":       s.TranslationCID,
		"verseListShowContent": s.VerseListShowContent,
	}
	if s.VerseListActiveID != nil {
		m["verseListActiveId"] = *s.VerseListActiveID
	} else {
		m["verseListActiveId"] = nil
	}

	return m
}

func mapToStudy(m map[string]any) StudySettings {
	columns := asStringSlice(m, "columns")
	if columns == nil {
		columns = []string{}
	}

	var activeID *string
	if v, ok := m["verseListActiveId"]; ok && v != nil {
		if s, ok := v.(string); ok {
			activeID = &s
		}
	}

	return StudySettings{
		ColumnCount:          asInt(m, "columnCount", 1),
		Columns:              columns,
		BookID:               asString(m, "bookId", "gen"),
		Chapter:              asInt(m, "chapter", 15),
		TranslationID:        asString(m, "translationId", "asv"),
		TranslationBID:       asString(m, "translationBId", "asv"),
		TranslationCID:       asString(m, "translationCId", "asv"),
		VerseListShowContent: asBool(m, "verseListShowContent", true),
		VerseListActiveID:    activeID,
	}
}

func stringSliceToAny(in []string) []any {
	out := make([]any, len(in))
	for i, s := range in {
		out[i] = s
	}

	return out
}
