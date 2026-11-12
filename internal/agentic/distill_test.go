//go:build agentic

package agentic

import "testing"

func TestRejectDistillDone(t *testing.T) {
	if rejectDistillDone(1, 0, "done", "every word") == "" {
		t.Fatal("expected reject with no upserts")
	}
	if rejectDistillDone(5, 3, "12 headwords beginning with A", "every word in the dictionary") == "" {
		t.Fatal("expected reject partial summary")
	}
	if rejectDistillDone(5, 3, "batch complete for chunk", "only letter A") != "" {
		t.Fatal("unexpected reject for scoped job")
	}
}

func TestFullDictionaryJob(t *testing.T) {
	if !fullDictionaryJob("Distill every headword") {
		t.Fatal("expected every")
	}
}

func TestMinUpsertsPerChunkConstant(t *testing.T) {
	if minUpsertsPerChunk < 3 {
		t.Fatal("expected meaningful minimum per chunk")
	}
}
