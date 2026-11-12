//go:build agentic

package agentic

import "testing"

func TestSlugify(t *testing.T) {
	if got := slugify("Ex. Of Faith"); got != "ex-of-faith" {
		t.Fatalf("got %q", got)
	}
}
