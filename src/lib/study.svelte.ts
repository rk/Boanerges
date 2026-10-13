import { Service, type StudySettings } from "../mygo";
import {
  columnTargetSlot,
  firstNonBibleSlot,
  normalizeColumns,
} from "./studyLayout";
import { columnDescriptor } from "./columns/catalog";

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
  const columnCount = (settings.columnCount as 1 | 2 | 3) || 1;
  study.columnCount = columnCount;
  study.columns = normalizeColumns(columnCount, (settings.columns ?? []) as ColumnType[]);
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
  const columns = normalizeColumns(study.columnCount, study.columns);
  study.columns = columns;
  return {
    columnCount: study.columnCount,
    columns,
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
  study.columns = normalizeColumns(count, study.columns);
  schedulePersistStudy();
}

export function setScrollSync(enabled: boolean): void {
  study.scrollSync = enabled;
}

export function setColumnContent(slotIndex: number, type: ColumnType): void {
  const columns = normalizeColumns(study.columnCount, [...study.columns]);
  columns[slotIndex] = type;
  study.columns = columns;
  schedulePersistStudy();
}

function resolveSlotIndex(type: ColumnType): number {
  const descriptor = columnDescriptor(type);
  if (descriptor.pickSlot === "firstNonBible") {
    const slot = firstNonBibleSlot(study.columns);
    if (slot !== null) {
      return slot;
    }
  }
  return columnTargetSlot(type, study.columnCount, study.columns);
}

/** Open or focus a study column (menu actions). */
export function openColumn(type: ColumnType): void {
  if (type === "search" && study.columns.includes("search")) {
    return;
  }
  if (study.columnCount === 1) {
    setColumnCount(2);
  }
  const slotIndex = resolveSlotIndex(type);
  if (study.columns[slotIndex] !== type) {
    setColumnContent(slotIndex, type);
  }
}

export function setSecondaryTranslation(slotIndex: number, translationId: string): void {
  if (slotIndex === 0) {
    study.translationBId = translationId;
  } else {
    study.translationCId = translationId;
  }
  schedulePersistStudy();
}

export function toggleColumn(type: ColumnType): void {
  openColumn(type);
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
