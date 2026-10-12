package bible

import "strings"

// ponytail: alias map only; full canon + BibleCanon lives in issue 3.

var canonOSIS = map[string]struct{}{
	"gen": {}, "exo": {}, "lev": {}, "num": {}, "deu": {}, "jos": {}, "jdg": {}, "rut": {},
	"1sa": {}, "2sa": {}, "1ki": {}, "2ki": {}, "1ch": {}, "2ch": {}, "ezr": {}, "neh": {},
	"est": {}, "job": {}, "psa": {}, "pro": {}, "ecc": {}, "sng": {}, "isa": {}, "jer": {},
	"lam": {}, "ezk": {}, "dan": {}, "hos": {}, "jol": {}, "amo": {}, "oba": {}, "jon": {},
	"mic": {}, "nam": {}, "hab": {}, "zep": {}, "hag": {}, "zec": {}, "mal": {}, "mat": {},
	"mrk": {}, "luk": {}, "jhn": {}, "act": {}, "rom": {}, "1co": {}, "2co": {}, "gal": {},
	"eph": {}, "php": {}, "col": {}, "1th": {}, "2th": {}, "1ti": {}, "2ti": {}, "tit": {},
	"phm": {}, "heb": {}, "jas": {}, "1pe": {}, "2pe": {}, "1jn": {}, "2jn": {}, "3jn": {},
	"jud": {}, "rev": {},
}

var osisAliases = map[string]string{
	"exod": "exo", "deut": "deu", "josh": "jos", "judg": "jdg", "1sam": "1sa", "2sam": "2sa",
	"1kgs": "1ki", "2kgs": "2ki", "1chr": "1ch", "2chr": "2ch", "matt": "mat", "matthew": "mat",
	"mark": "mrk", "luke": "luk", "john": "jhn", "genesis": "gen",
}

// NormalizeBookID returns the canonical OSIS id when input is known.
func NormalizeBookID(bookID string) (string, bool) {
	key := normalizeKey(bookID)
	if _, ok := canonOSIS[key]; ok {
		return key, true
	}
	if canonical, ok := osisAliases[key]; ok {
		return canonical, true
	}

	return "", false
}

func normalizeKey(bookID string) string {
	replacer := strings.NewReplacer(" ", "", ".", "", "-", "")
	return strings.ToLower(replacer.Replace(bookID))
}
