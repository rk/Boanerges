import { StudyData, type ScribeVerse } from "../mygo";
import type { ScribeDraftEntry } from "./scribeParagraphBreaks";

export {
  applyPrimaryParagraphBreaks,
  effectiveParagraphStart,
  resetInterVerseParagraphBreaks,
  type ScribeDraftEntry,
} from "./scribeParagraphBreaks";

export const scribe = $state({
  saveStatus: "idle" as "idle" | "saving" | "saved",
});

let saveTimeout: ReturnType<typeof setTimeout> | null = null;
let hideTimeout: ReturnType<typeof setTimeout> | null = null;
let saveGeneration = 0;

export async function fetchScribeDraft(bookId: string, chapter: number): Promise<ScribeVerse[]> {
  return StudyData.scribeGet(bookId, chapter);
}

export function scheduleScribeSave(bookId: string, chapter: number, verses: ScribeVerse[]): void {
  scribe.saveStatus = "saving";

  if (saveTimeout) {
    clearTimeout(saveTimeout);
  }
  if (hideTimeout) {
    clearTimeout(hideTimeout);
  }

  const generation = ++saveGeneration;

  saveTimeout = setTimeout(() => {
    void StudyData.scribePut(bookId, chapter, verses)
      .then(() => {
        if (generation !== saveGeneration) {
          return;
        }

        scribe.saveStatus = "saved";

        hideTimeout = setTimeout(() => {
          if (generation === saveGeneration) {
            scribe.saveStatus = "idle";
          }
        }, 1500);
      })
      .catch(() => {
        if (generation === saveGeneration) {
          scribe.saveStatus = "idle";
        }
      });
  }, 500);
}

export function entriesFromScribeVerses(verses: ScribeVerse[]): Record<number, ScribeDraftEntry> {
  const entries: Record<number, ScribeDraftEntry> = {};

  for (const verse of verses) {
    entries[verse.verse] = {
      text: verse.text,
      ...(verse.paragraphStart !== undefined ? { paragraphStartOverride: verse.paragraphStart } : {}),
    };
  }

  return entries;
}

export function serializeScribeDraft(
  verseNumbers: number[],
  entries: Record<number, ScribeDraftEntry>,
): ScribeVerse[] {
  const verses: ScribeVerse[] = [];

  for (const verseNumber of verseNumbers) {
    const entry = entries[verseNumber];
    const text = entry?.text ?? "";

    if (text.trim() === "" && entry?.paragraphStartOverride === undefined) {
      continue;
    }

    const serialized: ScribeVerse = {
      verse: verseNumber,
      text,
    };

    if (entry?.paragraphStartOverride !== undefined) {
      serialized.paragraphStart = entry.paragraphStartOverride;
    }

    verses.push(serialized);
  }

  return verses;
}
