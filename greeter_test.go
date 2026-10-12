package main

import "testing"

func TestGreeterGreet(t *testing.T) {
	var g Greeter
	if got := g.Greet("Ada"); got != "Hello, Ada! This message comes from Go." {
		t.Fatalf("Greet() = %q", got)
	}
	if got := g.Greet(""); got != "Hello, stranger! This message comes from Go." {
		t.Fatalf("Greet(empty) = %q", got)
	}
}
