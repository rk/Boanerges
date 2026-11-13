//go:build agentic

package agentic

import "testing"

func TestTopicalCannotUseDictionaryTools(t *testing.T) {
	for _, tool := range []string{"read_source", "upsert_entry", "list_sources", "lookup_webster"} {
		if toolAllowed("from_topic", tool) {
			t.Fatalf("topical should not allow %s", tool)
		}
	}
	for _, tool := range []string{"search", "search_bible", "read_reference", "cross_references", "upsert_topic"} {
		if !toolAllowed("from_topic", tool) {
			t.Fatalf("topical should allow %s", tool)
		}
	}
	if toolAllowed("from_topic", "get_chapter") {
		t.Fatal("topical should not use get_chapter; use read_reference")
	}
}

func TestDistillCannotUseTopicalOnlyPaths(t *testing.T) {
	if !toolAllowed("distill", "read_source") {
		t.Fatal("distill needs read_source")
	}
}
