package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func testRepo(t *testing.T) *Repository {
	t.Helper()

	return &Repository{Dir: t.TempDir()}
}

func writeRawSettings(t *testing.T, dir string, data map[string]any) {
	t.Helper()

	path := filepath.Join(dir, "settings.json")
	bytes, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, bytes, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestReadabilityDefaults(t *testing.T) {
	repo := testRepo(t)
	got := NewReadabilityStore(repo).Get()

	if got.FontSize != 18 || got.Theme != "auto" || got.FontFamily != "serif" {
		t.Fatalf("defaults = %+v", got)
	}
}

func TestReadabilityUpdate(t *testing.T) {
	repo := testRepo(t)
	store := NewReadabilityStore(repo)

	updated, err := store.Update(ReadabilitySettings{
		FontSize:    20,
		LineHeight:  1.8,
		Theme:       "dark",
		FontFamily:  "sans-serif",
		JustifyText: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if updated.FontSize != 20 || updated.Theme != "dark" {
		t.Fatalf("updated = %+v", updated)
	}

	got := store.Get()
	if got.FontSize != 20 || got.Theme != "dark" || got.FontFamily != "sans-serif" {
		t.Fatalf("persisted = %+v", got)
	}
}

func TestReadabilitySepiaMigratesOnRead(t *testing.T) {
	repo := testRepo(t)
	writeRawSettings(t, repo.Dir, map[string]any{
		"readability": map[string]any{"theme": "sepia"},
	})

	got := NewReadabilityStore(repo).Get()
	if got.Theme != "caramellatte" {
		t.Fatalf("theme = %q", got.Theme)
	}
}

func TestReadabilityValidation(t *testing.T) {
	svc := NewService(testRepo(t))

	_, err := svc.UpdateReadability(ReadabilitySettings{
		FontSize:    10,
		LineHeight:  1.8,
		Theme:       "dark",
		FontFamily:  "sans-serif",
		JustifyText: true,
	})
	if err == nil {
		t.Fatal("expected fontSize validation error")
	}

	_, err = svc.UpdateReadability(ReadabilitySettings{
		FontSize:    18,
		LineHeight:  1.7,
		Theme:       "sepia",
		FontFamily:  "serif",
		JustifyText: true,
	})
	if err == nil {
		t.Fatal("expected sepia theme rejection")
	}
}

func TestReadabilityAcceptsDaisyUIThemes(t *testing.T) {
	svc := NewService(testRepo(t))

	for _, theme := range []string{"nord", "dracula", "caramellatte", "auto"} {
		got, err := svc.UpdateReadability(ReadabilitySettings{
			FontSize:    18,
			LineHeight:  1.7,
			Theme:       theme,
			FontFamily:  "serif",
			JustifyText: true,
		})
		if err != nil {
			t.Fatalf("theme %q: %v", theme, err)
		}
		if got.Theme != theme {
			t.Fatalf("theme %q got %q", theme, got.Theme)
		}
	}
}

func TestStudyDefaults(t *testing.T) {
	got := NewStudyStore(testRepo(t)).Get()
	if got.ColumnCount != 1 || got.BookID != "gen" || got.Chapter != 15 {
		t.Fatalf("defaults = %+v", got)
	}
	if len(got.Columns) != 0 {
		t.Fatalf("columns = %v", got.Columns)
	}
}

func TestStudyUpdateViaService(t *testing.T) {
	svc := NewService(testRepo(t))

	got, err := svc.UpdateStudy(StudySettings{
		ColumnCount:          2,
		Columns:              []string{ColumnBibleSecondary},
		BookID:               "mat",
		Chapter:              5,
		TranslationID:        "asv",
		TranslationBID:       "web",
		TranslationCID:       "asv",
		VerseListShowContent: true,
		VerseListActiveID:    nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.ColumnCount != 2 || got.BookID != "mat" || got.Chapter != 5 {
		t.Fatalf("got = %+v", got)
	}

	persisted := svc.GetStudy()
	if persisted.BookID != "mat" {
		t.Fatalf("persisted bookId = %q", persisted.BookID)
	}
}

func TestStudyLegacyActiveViewMigration(t *testing.T) {
	repo := testRepo(t)
	writeRawSettings(t, repo.Dir, map[string]any{
		"study": map[string]any{
			"activeView":     "comparison",
			"bookId":         "mat",
			"chapter":        5,
			"translationId":  "asv",
			"translationBId": "web",
		},
	})

	got := NewStudyStore(repo).Get()
	if got.ColumnCount != 2 || len(got.Columns) != 1 || got.Columns[0] != ColumnBibleSecondary {
		t.Fatalf("layout = %+v", got)
	}
	if got.BookID != "mat" {
		t.Fatalf("bookId = %q", got.BookID)
	}
}

func TestStudyInvalidColumnSanitizedOnStoreUpdate(t *testing.T) {
	store := NewStudyStore(testRepo(t))

	store.Update(StudySettings{
		ColumnCount:          2,
		Columns:              []string{"invalid-type"},
		BookID:               "gen",
		Chapter:              1,
		TranslationID:        "asv",
		TranslationBID:       "asv",
		TranslationCID:       "asv",
		VerseListShowContent: true,
	})

	got := store.Get()
	if got.ColumnCount != 2 || len(got.Columns) != 1 || got.Columns[0] != ColumnBibleSecondary {
		t.Fatalf("sanitized = %+v", got)
	}
}

func TestStudyColumnTypesViaService(t *testing.T) {
	svc := NewService(testRepo(t))
	base := StudySettings{
		ColumnCount:          2,
		BookID:               "gen",
		Chapter:              1,
		TranslationID:        "asv",
		TranslationBID:       "asv",
		TranslationCID:       "asv",
		VerseListShowContent: true,
	}

	for _, column := range []string{ColumnCrossReferences, ColumnDictionary, ColumnComparison, ColumnVerseList} {
		settings := base
		settings.Columns = []string{column}
		got, err := svc.UpdateStudy(settings)
		if err != nil {
			t.Fatalf("column %q: %v", column, err)
		}
		if len(got.Columns) != 1 || got.Columns[0] != column {
			t.Fatalf("column %q got %+v", column, got.Columns)
		}
	}
}

func TestStudyValidationRejectsColumnSlotMismatch(t *testing.T) {
	svc := NewService(testRepo(t))

	_, err := svc.UpdateStudy(StudySettings{
		ColumnCount:          2,
		Columns:              []string{},
		BookID:               "gen",
		Chapter:              15,
		TranslationID:        "kjv",
		TranslationBID:       "asv",
		TranslationCID:       "asv",
		VerseListShowContent: true,
	})
	if err == nil {
		t.Fatal("expected column slot validation error")
	}
}
