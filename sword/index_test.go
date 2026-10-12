package sword

import (
	"path/filepath"
	"testing"
)

func TestGenesisIndexEntry(t *testing.T) {
	dir := filepath.Join("..", "extras", "sword", "modules", "texts", "ztext", "asv")
	td, err := loadTestament(dir, "ot")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := td.verseIndex["gen.1.1"]; !ok {
		var sample []string
		for k := range td.verseIndex {
			if len(sample) < 5 {
				sample = append(sample, k)
			}
		}
		t.Fatalf("Gen.1.1 missing; index size %d sample %v", len(td.verseIndex), sample)
	}
}
