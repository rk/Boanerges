//go:build agentic

package agentic

import (
	"strings"
	"testing"
)

func TestRejectTopicalDoneWithoutUpsert(t *testing.T) {
	if rejectTopicalDone("from_reference", 5, 0, 2, "saved topics") == "" {
		t.Fatal("expected reject without upsert_topic")
	}
}

func TestRejectTopicalDonePrematureSummary(t *testing.T) {
	s := "Identified and reused existing topical topics. No new topics were needed."
	if rejectTopicalDone("from_reference", 3, 0, 0, s) == "" {
		t.Fatal("expected reject for planning-only summary")
	}
}

func TestRejectTopicalDoneAcceptsAfterUpsert(t *testing.T) {
	if rejectTopicalDone("from_reference", 8, 3, 3, "Wrote eternal-life, salvation, love-of-neighbor") != "" {
		t.Fatal("expected accept when upserts and subagents done")
	}
}

func TestRejectFromReferenceNeedsSubagent(t *testing.T) {
	if rejectTopicalDone("from_reference", 4, 1, 0, "upserted one topic") == "" {
		t.Fatal("from_reference should require at least one run_subagent")
	}
}

func TestRejectFakeUpsertSummary(t *testing.T) {
	s := "Identified topics. Upserted 'Standing Instruction' topic."
	reason := rejectTopicalDone("from_reference", 4, 0, 0, s)
	if reason == "" || !strings.Contains(reason, "never ran") {
		t.Fatalf("want fake-upsert reject, got %q", reason)
	}
}
