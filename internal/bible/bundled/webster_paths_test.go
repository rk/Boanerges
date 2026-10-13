package bundled_test

import (
	"testing"

	"boanerges/internal/bible/bundled"
)

func TestFindWebsterDictionaryJSONFromRepo(t *testing.T) {
	path, err := bundled.FindWebsterDictionaryJSON("")
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		t.Fatal("expected path")
	}
}
