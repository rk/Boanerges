package bible

import (
	"slices"
	"testing"
)

func TestNormalizeBookID(t *testing.T) {
	cases := map[string]string{
		"matt": "mat", "mark": "mrk", "Matthew": "mat", "1Sam": "1sa",
		"gen": "gen", "ezra": "ezr", "zech": "zec", "acts": "act", "phlm": "phm",
	}
	for in, want := range cases {
		got, ok := NormalizeBookID(in)
		if !ok || got != want {
			t.Fatalf("NormalizeBookID(%q) = %q, %v; want %q", in, got, ok, want)
		}
	}
}

func TestDisplayName(t *testing.T) {
	if DisplayName("ezr") != "Ezra" || DisplayName("zec") != "Zechariah" || DisplayName("mrk") != "Mark" {
		t.Fatalf("display names mismatch")
	}
}

func TestLookupValues(t *testing.T) {
	mat := LookupValues("mat")
	if !slices.Contains(mat, "matt") || !slices.Contains(mat, "mat") {
		t.Fatalf("mat lookup = %v", mat)
	}
	mrk := LookupValues("mrk")
	if !slices.Contains(mrk, "mark") || !slices.Contains(mrk, "mrk") {
		t.Fatalf("mrk lookup = %v", mrk)
	}
}
