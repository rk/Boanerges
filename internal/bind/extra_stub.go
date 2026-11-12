//go:build !agentic

package bind

import "boanerges/internal/bibleapp"

func Extra(func(...any), *bibleapp.BibleService, *bibleapp.DictionaryService) {}
