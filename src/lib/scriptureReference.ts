import { canonBookName } from "./canonBookNames";
import type { Book } from "../mygo";

export type ScriptureReference = {
  bookId: string;
  chapter: number;
  verse: number;
};

/** Compact citation for search hits (e.g. Gen 1:1). */
export function formatScriptureReference(bookId: string, chapter: number, verse: number): string {
  return `${refAbbrev(bookId)} ${chapter}:${verse}`;
}

const REF_ABBREV: Record<string, string> = {
  gen: "Gen",
  exo: "Exod",
  lev: "Lev",
  num: "Num",
  deu: "Deut",
  jos: "Josh",
  jdg: "Judg",
  rut: "Ruth",
  "1sa": "1 Sam",
  "2sa": "2 Sam",
  "1ki": "1 Kgs",
  "2ki": "2 Kgs",
  "1ch": "1 Chr",
  "2ch": "2 Chr",
  ezr: "Ezra",
  neh: "Neh",
  est: "Esth",
  job: "Job",
  psa: "Ps",
  pro: "Prov",
  ecc: "Eccl",
  sng: "Song",
  isa: "Isa",
  jer: "Jer",
  lam: "Lam",
  ezk: "Ezek",
  dan: "Dan",
  hos: "Hos",
  jol: "Joel",
  amo: "Amos",
  oba: "Obad",
  jon: "Jonah",
  mic: "Mic",
  nam: "Nah",
  hab: "Hab",
  zep: "Zeph",
  hag: "Hag",
  zec: "Zech",
  mal: "Mal",
  mat: "Matt",
  mrk: "Mark",
  luk: "Luke",
  jhn: "John",
  act: "Acts",
  rom: "Rom",
  "1co": "1 Cor",
  "2co": "2 Cor",
  gal: "Gal",
  eph: "Eph",
  php: "Phil",
  col: "Col",
  "1th": "1 Thess",
  "2th": "2 Thess",
  "1ti": "1 Tim",
  "2ti": "2 Tim",
  tit: "Titus",
  phm: "Phlm",
  heb: "Heb",
  jas: "Jas",
  "1pe": "1 Pet",
  "2pe": "2 Pet",
  "1jn": "1 John",
  "2jn": "2 John",
  "3jn": "3 John",
  jud: "Jude",
  rev: "Rev",
};

function refAbbrev(bookId: string): string {
  const key = bookId.toLowerCase();
  return REF_ABBREV[key] ?? canonBookName(key).slice(0, 8);
}

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
