import { canonBookName } from "./canonBookNames";
import type { Book } from "../mygo";

export type ScriptureReference = {
  bookId: string;
  chapter: number;
  verse: number;
};

export function parseScriptureReference(input: string, books: Book[]): ScriptureReference | null {
  const trimmed = input.trim();
  if (trimmed === "") {
    return null;
  }

  const match = trimmed.match(/^(.+)\s+(\d+)\s*:\s*(\d+)\s*$/);
  if (match === null) {
    return null;
  }

  const bookPart = normalizeReferenceToken(match[1] ?? "");
  const chapter = Number.parseInt(match[2] ?? "", 10);
  const verse = Number.parseInt(match[3] ?? "", 10);
  if (chapter <= 0 || verse <= 0) {
    return null;
  }

  const book = resolveBook(bookPart, books);
  if (book === null) {
    return null;
  }

  return { bookId: book.id, chapter, verse };
}

function normalizeReferenceToken(value: string): string {
  return value.trim().toLowerCase().replace(/\./g, "").replace(/\s+/g, " ");
}

function resolveBook(bookPart: string, books: Book[]): Book | null {
  const normalized = normalizeReferenceToken(bookPart);
  const candidates = [...books].sort(
    (left, right) => bookMatchKeys(right).length - bookMatchKeys(left).length,
  );

  for (const book of candidates) {
    if (bookMatchKeys(book).some((key) => key === normalized)) {
      return book;
    }
  }

  return null;
}

function bookMatchKeys(book: Book): string[] {
  return [
    normalizeReferenceToken(book.name),
    normalizeReferenceToken(book.id),
    normalizeReferenceToken(book.abbrev),
    normalizeReferenceToken(canonBookName(book.id)),
  ];
}
