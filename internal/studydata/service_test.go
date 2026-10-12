package studydata_test

import (
	"testing"

	"boanerges/internal/studydata"
)

func TestNotesRoundTrip(t *testing.T) {
	svc, err := studydata.NewService()
	if err != nil {
		t.Skip("mygo user data path unavailable:", err)
	}

	if err := svc.NotesPut("gen", 1, "hello notes"); err != nil {
		t.Fatal(err)
	}
	got, err := svc.NotesGet("gen", 1)
	if err != nil || got != "hello notes" {
		t.Fatalf("notes = %q err=%v", got, err)
	}
}
