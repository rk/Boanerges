package openbible

import (
	"regexp"
	"strings"
	"sync"

	"boanerges/internal/bible"
)

var openBibleBooks = map[string]string{
	"Gen": "gen", "Exod": "exo", "Lev": "lev", "Num": "num", "Deut": "deu",
	"Josh": "jos", "Judg": "jdg", "Ruth": "rut", "1Sam": "1sa", "2Sam": "2sa",
	"1Kgs": "1ki", "2Kgs": "2ki", "1Chr": "1ch", "2Chr": "2ch", "Ezra": "ezr",
	"Neh": "neh", "Esth": "est", "Job": "job", "Ps": "psa", "Prov": "pro",
	"Eccl": "ecc", "Song": "sng", "Isa": "isa", "Jer": "jer", "Lam": "lam",
	"Ezek": "ezk", "Dan": "dan", "Hos": "hos", "Joel": "jol", "Amos": "amo",
	"Obad": "oba", "Jonah": "jon", "Mic": "mic", "Nah": "nam", "Hab": "hab",
	"Zeph": "zep", "Hag": "hag", "Zech": "zec", "Mal": "mal", "Matt": "mat",
	"Mark": "mrk", "Luke": "luk", "John": "jhn", "Acts": "act", "Rom": "rom",
	"1Cor": "1co", "2Cor": "2co", "Gal": "gal", "Eph": "eph", "Phil": "php",
	"Col": "col", "1Thess": "1th", "2Thess": "2th", "1Tim": "1ti", "2Tim": "2ti",
	"Titus": "tit", "Phlm": "phm", "Heb": "heb", "Jas": "jas", "1Pet": "1pe",
	"2Pet": "2pe", "1John": "1jn", "2John": "2jn", "3John": "3jn", "Jude": "jud",
	"Rev": "rev",
}

type Ref struct {
	BookID  string `json:"book_id"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
}

var (
	lookupOnce sync.Once
	lookup     map[int]Ref
	osisIndex  map[string]int
)

func FromVerseID(id int) *Ref {
	lookupOnce.Do(buildLookup)
	if ref, ok := lookup[id]; ok {
		r := ref

		return &r
	}

	return nil
}

func All() map[int]Ref {
	lookupOnce.Do(buildLookup)
	out := make(map[int]Ref, len(lookup))
	for k, v := range lookup {
		out[k] = v
	}

	return out
}

var refToken = regexp.MustCompile(`^(\d?[A-Za-z]+)\.(\d+)\.(\d+)$`)

type VerseIDRange struct {
	Start int
	End   *int
}

func VerseIDsFromReference(reference string) *VerseIDRange {
	reference = strings.TrimSpace(reference)
	if reference == "" {
		return nil
	}

	parts := strings.SplitN(reference, "-", 2)
	start := verseIDFromToken(parts[0])
	if start == 0 {
		return nil
	}
	if len(parts) == 1 {
		return &VerseIDRange{Start: start}
	}
	end := verseIDFromToken(parts[1])
	if end == 0 {
		return &VerseIDRange{Start: start}
	}

	return &VerseIDRange{Start: start, End: &end}
}

func verseIDFromToken(token string) int {
	m := refToken.FindStringSubmatch(strings.TrimSpace(token))
	if len(m) != 4 {
		return 0
	}
	osis, ok := openBibleBooks[m[1]]
	if !ok {
		return 0
	}
	lookupOnce.Do(buildLookup)
	key := osis + "." + m[2] + "." + m[3]

	return osisIndex[key]
}

func buildLookup() {
	lookup = map[int]Ref{}
	osisIndex = map[string]int{}
	id := 1
	for _, book := range bible.Books() {
		for chIdx, verseCount := range book.Chapters {
			chapter := chIdx + 1
			for verse := 1; verse <= verseCount; verse++ {
				ref := Ref{BookID: book.Osis, Chapter: chapter, Verse: verse}
				lookup[id] = ref
				osisIndex[book.Osis+"."+itoa(chapter)+"."+itoa(verse)] = id
				id++
			}
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b [12]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	return string(b[i:])
}
