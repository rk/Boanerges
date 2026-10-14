package settings

import "fmt"

func ValidateReadability(s ReadabilitySettings) error {
	if s.FontSize < 12 || s.FontSize > 32 {
		return validationError("fontSize must be between 12 and 32")
	}
	if s.LineHeight < 0.8 || s.LineHeight > 2 {
		return validationError("lineHeight must be between 0.8 and 2")
	}
	if !isReaderTheme(s.Theme) {
		return validationError("invalid theme")
	}
	if s.FontFamily != "sans-serif" && s.FontFamily != "serif" {
		return validationError("invalid fontFamily")
	}

	return nil
}

func ValidateStudy(s StudySettings) error {
	if s.ColumnCount != 1 && s.ColumnCount != 2 && s.ColumnCount != 3 {
		return validationError("invalid columnCount")
	}
	if s.Chapter < 1 || s.Chapter > 150 {
		return validationError("invalid chapter")
	}
	if len(s.BookID) == 0 || len(s.BookID) > 10 {
		return validationError("invalid bookId")
	}
	for _, id := range []string{s.TranslationID, s.TranslationBID, s.TranslationCID} {
		if len(id) == 0 || len(id) > 10 {
			return validationError("invalid translation id")
		}
	}
	if s.VerseListActiveID != nil && len(*s.VerseListActiveID) > 36 {
		return validationError("invalid verseListActiveId")
	}

	needed := max(0, s.ColumnCount-1)
	if len(s.Columns) != needed {
		return validationError("Column slots must match column count.")
	}

	nonBible := map[string]int{}
	for i, column := range s.Columns {
		if !allowedColumnTypes[column] {
			return validationError(fmt.Sprintf("invalid column type at index %d", i))
		}
		if columnAllowsDuplicate(column) {
			continue
		}
		if _, dup := nonBible[column]; dup {
			return validationError(fmt.Sprintf("Duplicate column type is not allowed at index %d", i))
		}
		nonBible[column] = i
	}

	bibleSlots := []int{}
	for i, column := range s.Columns {
		if column == ColumnBibleSecondary {
			bibleSlots = append(bibleSlots, i)
		}
	}
	if len(bibleSlots) == 2 {
		translations := []string{s.TranslationBID, s.TranslationCID}
		if translations[bibleSlots[0]] == translations[bibleSlots[1]] {
			return validationError("Secondary Bible columns must use different translations.")
		}
	}

	return nil
}
