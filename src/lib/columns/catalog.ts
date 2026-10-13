import type { ColumnType } from "../study.svelte.ts";

export type ColumnContentType = ColumnType;

export type ColumnPickSlot = "firstNonBible";

export type ColumnDescriptor = {
  type: ColumnContentType;
  label: string;
  allowDuplicate?: boolean;
  menuId?: string;
  pickSlot?: ColumnPickSlot;
};

export const COLUMN_CATALOG: ColumnDescriptor[] = [
  { type: "bible-secondary", label: "Translation", allowDuplicate: true },
  { type: "notes", label: "Notes", menuId: "study.notes" },
  { type: "scribe", label: "Scribe", menuId: "study.scribe" },
  { type: "search", label: "Search", menuId: "study.search", pickSlot: "firstNonBible" },
  { type: "cross-references", label: "Cross References", menuId: "study.cross-references" },
  { type: "comparison", label: "Comparison", menuId: "study.comparison" },
  { type: "verse-list", label: "Verse List", menuId: "study.verse-list" },
  { type: "dictionary", label: "Dictionary", menuId: "study.dictionary" },
];

const catalogByType = new Map(COLUMN_CATALOG.map((d) => [d.type, d]));

export const COLUMN_CONTENT_LABELS = Object.fromEntries(
  COLUMN_CATALOG.map((d) => [d.type, d.label]),
) as Record<ColumnContentType, string>;

export const COLUMN_CONTENT_TYPES: ColumnContentType[] = COLUMN_CATALOG.map((d) => d.type);

export function columnAllowsDuplicate(type: ColumnContentType): boolean {
  return catalogByType.get(type)?.allowDuplicate ?? false;
}

export function columnDescriptor(type: ColumnContentType): ColumnDescriptor {
  const descriptor = catalogByType.get(type);
  if (!descriptor) {
    throw new Error(`Unknown column type: ${type}`);
  }
  return descriptor;
}
