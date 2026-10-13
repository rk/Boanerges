package studydata_test

import (
	"os"
	"path/filepath"
	"testing"

	"boanerges/internal/studydata"
)

func TestLegacyNotesPathLayout(t *testing.T) {
	root := t.TempDir()
	svc := studydata.NewServiceForTest(root)

	if err := svc.NotesPut("gen", 1, "legacy note"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "notes", "gen", "1.txt")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(raw) != "legacy note" {
		t.Fatalf("notes file = %q", raw)
	}
}

func TestLegacyScribeJSONLayout(t *testing.T) {
	root := t.TempDir()
	svc := studydata.NewServiceForTest(root)

	verses := []studydata.ScribeVerse{{Verse: 1, Text: "In the beginning", ParagraphStart: true}}
	if err := svc.ScribePut("gen", 1, verses); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "scribe", "gen", "1.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
	got, err := svc.ScribeGet("gen", 1)
	if err != nil || len(got) != 1 || got[0].Text != "In the beginning" {
		t.Fatalf("scribe = %+v err=%v", got, err)
	}
}

func TestLegacyVerseListIndexLayout(t *testing.T) {
	root := t.TempDir()
	svc := studydata.NewServiceForTest(root)

	list, err := svc.VerseListSave("My list", []studydata.VerseListEntry{
		{BookID: "gen", Chapter: 1, Verse: 1},
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "verse-lists", list.ID+".json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "verse-lists", "index.json")); err != nil {
		t.Fatal(err)
	}
	idx, err := svc.VerseListIndex()
	if err != nil || len(idx) != 1 || idx[0].Title != "My list" {
		t.Fatalf("index = %+v err=%v", idx, err)
	}
}

func TestVerseListRejectsInvalidEntry(t *testing.T) {
	svc := studydata.NewServiceForTest(t.TempDir())
	_, err := svc.VerseListSave("Bad", []studydata.VerseListEntry{
		{BookID: "INVALID", Chapter: 1, Verse: 1},
	}, "")
	if err == nil {
		t.Fatal("expected invalid book id error")
	}
}

func TestVerseListDelete(t *testing.T) {
	root := t.TempDir()
	svc := studydata.NewServiceForTest(root)
	list, err := svc.VerseListSave("Gone", []studydata.VerseListEntry{{BookID: "gen", Chapter: 1, Verse: 1}}, "")
	if err != nil {
		t.Fatal(err)
	}
	if err := svc.VerseListDelete(list.ID); err != nil {
		t.Fatal(err)
	}
	idx, err := svc.VerseListIndex()
	if err != nil || len(idx) != 0 {
		t.Fatalf("index after delete = %+v err=%v", idx, err)
	}
}
