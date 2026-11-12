import { describe, expect, it } from "vitest";
import { verseListEntryKey } from "./study.svelte";

describe("verseListEntryKey", () => {
  it("dedupes identical references", () => {
    const a = { bookId: "gen", chapter: 1, verse: 1 };
    const b = { bookId: "gen", chapter: 1, verse: 1 };
    expect(verseListEntryKey(a)).toBe(verseListEntryKey(b));
  });
});
