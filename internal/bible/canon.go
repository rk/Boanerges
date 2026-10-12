package bible

import (
	_ "embed"
	"encoding/json"
)

//go:embed canon.json
var canonJSON []byte

// CanonBook is one book in the Protestant canon.
type CanonBook struct {
	Osis     string `json:"osis"`
	Chapters []int  `json:"chapters"`
}

var canonBooks []CanonBook

func init() {
	if err := json.Unmarshal(canonJSON, &canonBooks); err != nil {
		panic(err)
	}
}

// Books returns the embedded canon (osis id + chapter verse counts).
func Books() []CanonBook {
	out := make([]CanonBook, len(canonBooks))
	copy(out, canonBooks)

	return out
}

func canonOSISSet() map[string]struct{} {
	set := make(map[string]struct{}, len(canonBooks))
	for _, book := range canonBooks {
		set[book.Osis] = struct{}{}
	}

	return set
}
