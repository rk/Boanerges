//go:build agentic

package agentic

import "testing"

func TestReadSourceResolveName(t *testing.T) {
	a := sourceChunkArgs{SourceFile: "chambers-1-of-4-A-D.txt", Offset: 0, Limit: 8000}
	name, err := a.resolveName()
	if err != nil || name != "chambers-1-of-4-A-D.txt" {
		t.Fatalf("resolve: %v %q", err, name)
	}
}

func TestValidateRejectsInstructionEcho(t *testing.T) {
	inst := "Distill every headword from every source"
	e := DictEntry{Headword: "Distill", Modern: inst, Senses: []DictSense{{Gloss: "x"}}}
	if err := validateDictEntry(e, inst); err == nil {
		t.Fatal("expected reject")
	}
}
