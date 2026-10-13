import { Service, StudyData, type StudySettings } from "../mygo";
import { normalizeDictionaryWord } from "./normalizeDictionaryWord";
import { formatScriptureReference } from "./scriptureReference";
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
  selectedVerse: null as number | null,
  /** Verse cross-ref panel should load (displayed passage, not stale study location). */
  crossRefsTarget: null as { bookId: string; chapter: number; verse: number } | null,
  crossRefsTick: 0,
  dictionaryQuery: null as string | null,
});

let hydrated = false;
let persistTimer: ReturnType<typeof setTimeout> | null = null;

/** Apply persisted study settings to reactive state (safe to call after flyout saves). */
export function applyStudySettings(settings: StudySettings): void {
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
}

export function hydrateStudy(settings: StudySettings): void {
  if (hydrated) return;
  applyStudySettings(settings);
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
  study.selectedVerse = null;
  schedulePersistStudy();
}

export function setChapter(chapter: number): void {
  study.chapter = chapter;
  study.selectedVerse = null;
  schedulePersistStudy();
}

export function setLocation(bookId: string, chapter: number): void {
  study.bookId = bookId;
  study.chapter = chapter;
  study.selectedVerse = null;
  schedulePersistStudy();
}

export function setActiveVerse(verse: number): void {
  study.selectedVerse = verse;
}

export function navigateToVerse(bookId: string, chapter: number, verse: number): void {
  study.bookId = bookId;
  study.chapter = chapter;
  study.selectedVerse = verse;
  schedulePersistStudy();
  requestAnimationFrame(() => {
    document.getElementById(`verse-${verse}`)?.scrollIntoView({ block: "center", behavior: "smooth" });
  });
}

export async function addVerseToList(bookId: string, chapter: number, verse: number): Promise<void> {
  let activeId = study.verseListActiveId;
  let entries: { bookId: string; chapter: number; verse: number }[] = [];
  let title = "My list";

  if (activeId) {
    const list = await StudyData.verseListGet(activeId);
    entries = [...list.entries];
    title = list.title;
  }

  entries.push({ bookId, chapter, verse });
  const saved = await StudyData.verseListSave(title, entries, activeId ?? "");
  study.verseListActiveId = saved.id;
  schedulePersistStudy();
}

export function openCrossReferencesForVerse(bookId: string, chapter: number, verse: number): void {
  if (study.bookId !== bookId || study.chapter !== chapter) {
    study.bookId = bookId;
    study.chapter = chapter;
    schedulePersistStudy();
  }
  study.selectedVerse = verse;
  study.crossRefsTarget = { bookId, chapter, verse };
  study.crossRefsTick += 1;
  openColumn("cross-references");
}

export function openDictionaryForWord(word: string): void {
  const normalized = normalizeDictionaryWord(word);
  if (normalized.length < 2) {
    return;
  }
  study.dictionaryQuery = normalized;
  openColumn("dictionary");
}

export async function copyVerseReference(bookId: string, chapter: number, verse: number): Promise<void> {
  const text = formatScriptureReference(bookId, chapter, verse);
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    // ponytail: clipboard may fail in webview without permission
  }
}
