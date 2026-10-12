package settings

// Study column type ids match legacy App\Enums\StudyColumnType.
const (
	ColumnBibleSecondary  = "bible-secondary"
	ColumnNotes           = "notes"
	ColumnScribe          = "scribe"
	ColumnSearch          = "search"
	ColumnCrossReferences = "cross-references"
	ColumnComparison      = "comparison"
	ColumnVerseList       = "verse-list"
	ColumnDictionary      = "dictionary"
)

var allowedColumnTypes = map[string]bool{
	ColumnBibleSecondary:  true,
	ColumnNotes:           true,
	ColumnScribe:          true,
	ColumnSearch:          true,
	ColumnCrossReferences: true,
	ColumnComparison:      true,
	ColumnVerseList:       true,
	ColumnDictionary:      true,
}

func columnAllowsDuplicate(column string) bool {
	return column == ColumnBibleSecondary
}
