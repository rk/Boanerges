package dictionary

import (
	"database/sql"
	"encoding/json"
	"strings"
	"unicode"
)

const MetaKey = "webster1828_dictionary_v1"

var posLabels = map[string]string{
	"n.": "noun", "v.": "verb", "a.": "adjective", "adv.": "adverb",
	"prep.": "preposition", "pron.": "pronoun", "conj.": "conjunction",
	"interj.": "interjection", "p.": "participle",
}

type Service struct {
	db *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) IsImported() bool {
	var completed sql.NullString
	if err := s.db.QueryRow(`SELECT completed_at FROM import_meta WHERE key = ?`, MetaKey).Scan(&completed); err != nil || !completed.Valid {
		return false
	}
	var n int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM dictionary_entries`).Scan(&n); err != nil || n == 0 {
		return false
	}
	return true
}

type Variant struct {
	PartOfSpeech         *string  `json:"partOfSpeech"`
	PartOfSpeechExpanded string   `json:"partOfSpeechExpanded"`
	Definitions          []string `json:"definitions"`
}

type LookupResult struct {
	Found    bool      `json:"found"`
	Word     string    `json:"word"`
	Variants []Variant `json:"variants,omitempty"`
	Message  string    `json:"message,omitempty"`
}

func (s *Service) Lookup(word string) LookupResult {
	key := normalizeWordKey(word)
	if key == "" {
		return LookupResult{Found: false, Word: word, Message: "Enter a valid word to look up."}
	}
	if !s.IsImported() {
		return LookupResult{Found: false, Word: word, Message: "Dictionary is still being set up."}
	}

	rows, err := s.db.Query(`
		SELECT word, part_of_speech, definitions FROM dictionary_entries
		WHERE word_key = ? ORDER BY id`, key)
	if err != nil {
		return LookupResult{Found: false, Word: word, Message: err.Error()}
	}
	defer rows.Close()

	var variants []Variant
	var displayWord string
	for rows.Next() {
		var w string
		var pos sql.NullString
		var defsJSON string
		if err := rows.Scan(&w, &pos, &defsJSON); err != nil {
			return LookupResult{Found: false, Word: word, Message: err.Error()}
		}
		if displayWord == "" {
			displayWord = w
		}
		var posPtr *string
		expanded := ""
		if pos.Valid {
			p := pos.String
			posPtr = &p
			expanded = expandPOS(p)
		}
		variants = append(variants, Variant{
			PartOfSpeech: posPtr, PartOfSpeechExpanded: expanded,
			Definitions: decodeDefs(defsJSON),
		})
	}
	if len(variants) == 0 {
		return LookupResult{Found: false, Word: word, Message: "Word '" + word + "' is not in this dictionary."}
	}

	return LookupResult{Found: true, Word: displayWord, Variants: variants}
}

func (s *Service) Suggest(prefix string, limit int) []string {
	if !s.IsImported() {
		return nil
	}
	key := normalizeWordKey(prefix)
	if len(key) < 2 {
		return nil
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 25 {
		limit = 25
	}

	rows, err := s.db.Query(`
		SELECT word, word_key FROM dictionary_entries
		WHERE word_key LIKE ? ORDER BY word_key LIMIT ?`, key+"%", limit*3)
	if err != nil {
		return nil
	}
	defer rows.Close()

	seen := map[string]struct{}{}
	var out []string
	for rows.Next() {
		var word, wordKey string
		if err := rows.Scan(&word, &wordKey); err != nil {
			return out
		}
		if _, ok := seen[wordKey]; ok {
			continue
		}
		seen[wordKey] = struct{}{}
		out = append(out, word)
		if len(out) >= limit {
			break
		}
	}

	return out
}

func normalizeWordKey(word string) string {
	word = strings.TrimSpace(word)
	if word == "" {
		return ""
	}
	token := strings.Fields(word)[0]
	token = strings.TrimFunc(token, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '\'' && r != '-'
	})

	return strings.ToLower(token)
}

func expandPOS(pos string) string {
	if pos == "" {
		return ""
	}
	if label, ok := posLabels[pos]; ok {
		return label
	}

	return pos
}

func decodeDefs(raw string) []string {
	var defs []string
	if err := json.Unmarshal([]byte(raw), &defs); err != nil {
		return nil
	}

	return defs
}
