export type CrossReferenceHit = {
  rank: number;
  bookId: string;
  bookName: string;
  chapter: number;
  verse: number;
  endVerse: number | null;
};

function readField(obj: Record<string, unknown>, ...keys: string[]): unknown {
  for (const key of keys) {
    if (key in obj) {
      return obj[key];
    }
  }
  return undefined;
}

function normalizeHit(raw: unknown, index: number): CrossReferenceHit | null {
  if (!raw || typeof raw !== "object") {
    return null;
  }
  const o = raw as Record<string, unknown>;
  const bookId = String(readField(o, "bookId", "BookID") ?? "");
  if (!bookId) {
    return null;
  }
  const bookName = String(readField(o, "bookName", "BookName") ?? bookId);
  const chapter = Number(readField(o, "chapter", "Chapter"));
  const verse = Number(readField(o, "verse", "Verse"));
  const rankRaw = readField(o, "rank", "Rank");
  const rank = rankRaw === undefined || rankRaw === null ? index + 1 : Number(rankRaw);
  const endRaw = readField(o, "endVerse", "EndVerse");
  const endVerse = endRaw === undefined || endRaw === null ? null : Number(endRaw);

  if (!Number.isFinite(chapter) || !Number.isFinite(verse)) {
    return null;
  }

  return { rank, bookId, bookName, chapter, verse, endVerse };
}

/** Normalize mygo cross-ref payload (OpenBible ranks are not unique per verse). */
export function normalizeCrossReferences(resp: unknown): CrossReferenceHit[] {
  if (!resp || typeof resp !== "object") {
    return [];
  }
  const raw = readField(resp as Record<string, unknown>, "references", "References");
  if (!Array.isArray(raw)) {
    return [];
  }
  const out: CrossReferenceHit[] = [];
  for (let i = 0; i < raw.length; i++) {
    const hit = normalizeHit(raw[i], i);
    if (hit) {
      out.push(hit);
    }
  }
  return out;
}

export function crossReferenceKey(ref: CrossReferenceHit, index: number): string {
  return `${index}-${ref.bookId}-${ref.chapter}-${ref.verse}-${ref.rank}`;
}
