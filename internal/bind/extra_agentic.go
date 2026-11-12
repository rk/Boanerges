//go:build agentic

package bind

import (
	"boanerges/internal/agentic"
	"boanerges/internal/bibleapp"
)

func Extra(bind func(...any), bible *bibleapp.BibleService, dict *bibleapp.DictionaryService) {
	bind(agentic.NewAgentic(bible, dict))
}
