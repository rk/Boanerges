//go:build agentic

package agentic

import (
	"strings"
	"testing"
)

func TestSearchArgsRejectReferenceShape(t *testing.T) {
	a := searchArgs{BookID: "jhn", Chapter: 3, Verse: 16}
	if err := a.validate(); err == nil || !strings.Contains(err.Error(), "read_reference") {
		t.Fatalf("want read_reference hint, got %v", err)
	}
}

func TestSearchArgsRequireQuery(t *testing.T) {
	var empty searchArgs
	if err := empty.validate(); err == nil {
		t.Fatal("expected error")
	}
}

func TestNormalizeSearchAlias(t *testing.T) {
	if normalizeToolName("search_bible") != "search" {
		t.Fatal("search_bible should alias to search")
	}
}

func TestReadReferenceArgsVerseRange(t *testing.T) {
	a := readReferenceArgs{BookID: "gen", Chapter: 1, Verse: 1, VerseEnd: 3}
	start, end, err := a.verseRange()
	if err != nil || start != 1 || end != 3 {
		t.Fatalf("range = %d-%d err=%v", start, end, err)
	}
}
