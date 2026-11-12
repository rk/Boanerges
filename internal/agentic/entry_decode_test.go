//go:build agentic

package agentic

import (
	"encoding/json"
	"testing"
)

func TestDecodeDictEntryWrappedAndLoose(t *testing.T) {
	raw := json.RawMessage(`{"entry":{"headword":"Able","modern":"capable","senses":[{"gloss":"having power"}]}}`)
	e, err := decodeDictEntryJSON(raw)
	if err != nil || e.Headword != "Able" {
		t.Fatalf("wrap: %v %+v", err, e)
	}
	raw = json.RawMessage(`{"word":"Faith","definition":"trust"}`)
	e, err = decodeDictEntryJSON(raw)
	if err != nil || e.Headword != "Faith" || len(e.Senses) != 1 {
		t.Fatalf("loose: %v %+v", err, e)
	}
}

func TestDecodeFromToolLineWithoutArgs(t *testing.T) {
	line := `{"tool":"upsert_entry","headword":"Aaron","modern":"masc. name","senses":[{"gloss":"...","sourceFile":"a.txt","sourceSpan":"1-50"}]}`
	list, err := decodeDictEntries(nil, line)
	if err != nil || len(list) != 1 || list[0].Headword != "Aaron" {
		t.Fatalf("line: %v %+v", err, list)
	}
}
