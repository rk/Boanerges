package sword

import (
	"strings"

	"boanerges/internal/bible"
)

// BookInfo describes one bible book for index walks (KJV versification).
type BookInfo struct {
	Name        string
	Abbrev      string
	MaxChapter  int
	Verses      []int // verse count per chapter
	Testament   string
	CanonOSIS   string
}

// KVJBooks returns 66 books in canonical order with KJV verse counts per chapter.
func KVJBooks() []BookInfo {
	kjv := loadKJVVersification()
	out := make([]BookInfo, 0, len(bible.Books()))
	for i, book := range bible.Books() {
		testament := "ot"
		if i >= 39 {
			testament = "nt"
		}
		abbrev := swordAbbrev(book.Osis)
		verses := kjv[i]
		if len(verses) == 0 {
			verses = book.Chapters
		}
		out = append(out, BookInfo{
			Name:       bible.DisplayName(book.Osis),
			Abbrev:     abbrev,
			MaxChapter: len(verses),
			Verses:     verses,
			Testament:  testament,
			CanonOSIS:  book.Osis,
		})
	}

	return out
}

// swordAbbrev maps canon OSIS ids to common SWORD OSIS abbreviations (KJV).
func swordAbbrev(canon string) string {
	switch canon {
	case "exo":
		return "Exod"
	case "deu":
		return "Deut"
	case "jos":
		return "Josh"
	case "jdg":
		return "Judg"
	case "1sa":
		return "1Sam"
	case "2sa":
		return "2Sam"
	case "1ki":
		return "1Kgs"
	case "2ki":
		return "2Kgs"
	case "1ch":
		return "1Chr"
	case "2ch":
		return "2Chr"
	case "est":
		return "Esth"
	case "psa":
		return "Ps"
	case "pro":
		return "Prov"
	case "ecc":
		return "Eccl"
	case "sng":
		return "Song"
	case "mat":
		return "Matt"
	case "mrk":
		return "Mark"
	case "luk":
		return "Luke"
	case "jhn":
		return "John"
	case "act":
		return "Acts"
	case "1co":
		return "1Cor"
	case "2co":
		return "2Cor"
	case "1th":
		return "1Thess"
	case "2th":
		return "2Thess"
	case "1ti":
		return "1Tim"
	case "2ti":
		return "2Tim"
	case "phm":
		return "Phlm"
	case "1pe":
		return "1Pet"
	case "2pe":
		return "2Pet"
	case "1jn":
		return "1John"
	case "2jn":
		return "2John"
	case "3jn":
		return "3John"
	default:
		if len(canon) == 0 {
			return canon
		}
		return stringsTitle(canon)
	}
}

func stringsTitle(s string) string {
	if len(s) <= 3 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
