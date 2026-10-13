import { columnAllowsDuplicate, COLUMN_CONTENT_TYPES, columnDescriptor } from "./columns/catalog";
import type { Book } from "../mygo";
import type { ColumnType } from "./study.svelte.ts";

export function activeBibleTranslationIds(settings: {
  translationId: string;
  translationBId: string;
  translationCId: string;
  columns: ColumnType[];
}): string[] {
  const ids = [settings.translationId];
  for (const [slotIndex, column] of settings.columns.entries()) {
    if (column === "bible-secondary") {
      ids.push(slotIndex === 0 ? settings.translationBId : settings.translationCId);
    }
  }

  return [...new Set(ids)];
}

export function isBookAvailableInTranslations(
  bookId: string,
  translationIds: string[],
  booksByTranslation: ReadonlyMap<string, readonly Book[]>,
): boolean {
  return translationIds.every((translationId) =>
    booksByTranslation.get(translationId)?.some((book) => book.id === bookId),
  );
}

export function bibleColumnCount(settings: { columnCount: number; columns: ColumnType[] }): number {
  return 1 + settings.columns.filter((column) => column === "bible-secondary").length;
}

export function normalizeColumns(columnCount: 1 | 2 | 3, columns: ColumnType[]): ColumnType[] {
  const needed = Math.max(0, columnCount - 1);
  const normalized = columns.slice(0, needed);
  while (normalized.length < needed) {
    normalized.push(defaultColumnForSlot(normalized));
  }
  return normalized;
}

function defaultColumnForSlot(existing: ColumnType[]): ColumnType {
  for (const type of COLUMN_CONTENT_TYPES) {
    if (columnAllowsDuplicate(type) || !existing.includes(type)) {
      return type;
    }
  }
  return "bible-secondary";
}

export function columnTargetSlot(type: ColumnType, columnCount: 1 | 2 | 3, columns: ColumnType[]): number {
  for (let index = columns.length - 1; index >= 0; index--) {
    if (columns[index] === type) {
      return index;
    }
  }
  return Math.max(0, columnCount - 2);
}

export function firstNonBibleSlot(columns: ColumnType[]): number | null {
  const index = columns.findIndex((column) => column !== "bible-secondary");
  return index >= 0 ? index : null;
}

export function secondaryTranslationForSlot(
  slotIndex: number,
  translationBId: string,
  translationCId: string,
): string {
  return slotIndex === 0 ? translationBId : translationCId;
}

export function availableColumnOptions(
  slotIndex: number,
  columns: ColumnType[],
  translationBId: string,
  translationCId: string,
): ColumnType[] {
  const used = new Set(columns.filter((_, index) => index !== slotIndex));

  return COLUMN_CONTENT_TYPES.filter((type) => {
    if (!columnAllowsDuplicate(type) && used.has(type)) {
      return false;
    }

    if (type !== "bible-secondary") {
      return true;
    }

    const otherBibleSlots = columns
      .map((column, index) => ({ column, index }))
      .filter(({ column, index }) => column === "bible-secondary" && index !== slotIndex);

    if (otherBibleSlots.length === 0) {
      return true;
    }

    const otherSlot = otherBibleSlots[0];
    if (!otherSlot) {
      return true;
    }
    const otherIndex = otherSlot.index;
    const otherTranslation = otherIndex === 0 ? translationBId : translationCId;
    const thisTranslation = slotIndex === 0 ? translationBId : translationCId;

    return otherTranslation !== thisTranslation;
  });
}
