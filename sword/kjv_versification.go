package sword

import (
	_ "embed"
	"encoding/json"
)

//go:embed kjv_verses_in_chapter.json
var kjvVersesJSON []byte

func loadKJVVersification() [][]int {
	var raw [][]int
	if err := json.Unmarshal(kjvVersesJSON, &raw); err != nil {
		return nil
	}

	return raw
}
