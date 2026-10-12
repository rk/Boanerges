package bible

import (
	"strings"
	"sync"
)

var (
	canonOnce sync.Once
	canonSet  map[string]struct{}

	lookupOnce sync.Once
	lookupMap  map[string][]string
)

// NormalizeBookID returns the canonical OSIS id when input is known.
func NormalizeBookID(bookID string) (string, bool) {
	key := normalizeKey(bookID)
	canonOnce.Do(func() { canonSet = canonOSISSet() })

	if _, ok := canonSet[key]; ok {
		return key, true
	}
	if canonical, ok := osisAliases[key]; ok {
		return canonical, true
	}

	return "", false
}

// IsCanon reports whether bookID resolves to a canon book.
func IsCanon(bookID string) bool {
	_, ok := NormalizeBookID(bookID)

	return ok
}

// DisplayName returns the English book name for UI.
func DisplayName(bookID string) string {
	canonical, ok := NormalizeBookID(bookID)
	if !ok {
		return strings.ToUpper(bookID)
	}
	if name, ok := osisDisplayNames[canonical]; ok {
		return name
	}

	return strings.ToUpper(canonical)
}

// LookupValues returns DB match keys for a canonical or alias id.
func LookupValues(bookID string) []string {
	canonical, ok := NormalizeBookID(bookID)
	if !ok {
		return []string{normalizeKey(bookID)}
	}

	lookupOnce.Do(func() { lookupMap = buildLookupValues() })

	if values, ok := lookupMap[canonical]; ok {
		out := make([]string, len(values))
		copy(out, values)

		return out
	}

	return []string{canonical}
}

func buildLookupValues() map[string][]string {
	values := map[string][]string{}
	for _, book := range Books() {
		values[book.Osis] = []string{book.Osis}
	}
	for alias, canonical := range osisAliases {
		values[canonical] = append(values[canonical], alias)
	}
	for canonical, aliases := range values {
		seen := map[string]struct{}{}
		unique := make([]string, 0, len(aliases))
		for _, a := range aliases {
			if _, ok := seen[a]; ok {
				continue
			}
			seen[a] = struct{}{}
			unique = append(unique, a)
		}
		values[canonical] = unique
	}

	return values
}

func normalizeKey(bookID string) string {
	replacer := strings.NewReplacer(" ", "", ".", "", "-", "")
	return strings.ToLower(replacer.Replace(bookID))
}
