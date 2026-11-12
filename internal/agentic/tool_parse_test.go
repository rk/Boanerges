//go:build agentic

package agentic

import (
	"testing"
)

func TestParseToolTurnMarkdownAndAlias(t *testing.T) {
	reply := "Here:\n```json\n{\"action\":\"upsert_entry\",\"args\":{\"headword\":\"Aaron\",\"modern\":\"...\",\"senses\":[{\"gloss\":\"g\"}]}}\n```"
	tool, args, done, _, err := parseToolTurn(reply)
	if err != nil || done || tool != "upsert_entry" || len(args) == 0 {
		t.Fatalf("got tool=%q done=%v err=%v args=%s", tool, done, err, args)
	}
}

func TestParseImplicitUpsert(t *testing.T) {
	reply := `I will add: {"headword":"Faith","senses":[{"gloss":"trust"}]}`
	tool, _, done, _, err := parseToolTurn(reply)
	if err != nil || done || tool != "upsert_entry" {
		t.Fatalf("implicit: tool=%q err=%v", tool, err)
	}
}

func TestParseDoneNotConfusedWithEntry(t *testing.T) {
	reply := `{"done":true,"summary":"ok"}`
	_, _, done, _, err := parseToolTurn(reply)
	if err != nil || !done {
		t.Fatal("expected done")
	}
}
