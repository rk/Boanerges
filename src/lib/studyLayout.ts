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
