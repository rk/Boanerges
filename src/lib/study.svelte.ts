import { Service, type StudySettings } from "../mygo";

export type ColumnType =
  | "bible-secondary"
  | "notes"
  | "scribe"
  | "search"
  | "cross-references"
  | "comparison"
  | "verse-list"
  | "dictionary";

export const study = $state({
  columnCount: 1 as 1 | 2 | 3,
  columns: [] as ColumnType[],
  bookId: "gen",
  chapter: 1,
  translationId: "asv",
  translationBId: "asv",
  translationCId: "asv",
  verseListShowContent: true,
  verseListActiveId: null as string | null,
  scrollSync: false,
  settingsOpen: false,
  selectedVerse: 1,
});

let hydrated = false;
let persistTimer: ReturnType<typeof setTimeout> | null = null;

export function hydrateStudy(settings: StudySettings): void {
  if (hydrated) return;
  study.columnCount = (settings.columnCount as 1 | 2 | 3) || 1;
  study.columns = (settings.columns ?? []) as ColumnType[];
  study.bookId = settings.bookId || "gen";
  study.chapter = settings.chapter || 1;
  study.translationId = settings.translationId || "asv";
  study.translationBId = settings.translationBId || "asv";
  study.translationCId = settings.translationCId || "asv";
  study.verseListShowContent = settings.verseListShowContent ?? true;
  study.verseListActiveId = settings.verseListActiveId ?? null;
  hydrated = true;
}

function payload(): StudySettings {
  return {
    columnCount: study.columnCount,
    columns: study.columns,
    bookId: study.bookId,
    chapter: study.chapter,
    translationId: study.translationId,
    translationBId: study.translationBId,
    translationCId: study.translationCId,
    verseListShowContent: study.verseListShowContent,
    verseListActiveId: study.verseListActiveId,
  };
}

export function schedulePersistStudy(): void {
  if (persistTimer) clearTimeout(persistTimer);
  persistTimer = setTimeout(() => {
    void Service.updateStudy(payload());
  }, 250);
}

export function setColumnCount(count: 1 | 2 | 3): void {
  study.columnCount = count;
  if (study.columns.length > count - 1) {
    study.columns = study.columns.slice(0, count - 1);
  }
  schedulePersistStudy();
}

export function toggleColumn(type: ColumnType): void {
  const idx = study.columns.indexOf(type);
  if (idx >= 0) {
    study.columns = study.columns.filter((c) => c !== type);
  } else if (study.columns.length < study.columnCount - 1) {
    study.columns = [...study.columns, type];
  }
  schedulePersistStudy();
}

export function setBook(bookId: string): void {
  study.bookId = bookId;
  study.chapter = 1;
  schedulePersistStudy();
}

export function setChapter(chapter: number): void {
  study.chapter = chapter;
  schedulePersistStudy();
}

export function setLocation(bookId: string, chapter: number): void {
  study.bookId = bookId;
  study.chapter = chapter;
  schedulePersistStudy();
}
