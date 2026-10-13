import { describe, expect, it } from "vitest";
import { crossReferenceKey, normalizeCrossReferences } from "./crossReferences";

describe("normalizeCrossReferences", () => {
  it("reads camelCase and PascalCase fields", () => {
    const hits = normalizeCrossReferences({
      References: [
        { Rank: 1, BookID: "jhn", BookName: "John", Chapter: 1, Verse: 1, EndVerse: null },
      ],
    });
    expect(hits).toHaveLength(1);
    expect(hits[0]?.bookName).toBe("John");
  });

  it("keeps duplicate ranks as separate hits", () => {
    const hits = normalizeCrossReferences({
      references: [
        { rank: 65, bookId: "a", bookName: "A", chapter: 1, verse: 1, endVerse: null },
        { rank: 65, bookId: "b", bookName: "B", chapter: 2, verse: 2, endVerse: null },
      ],
    });
    expect(hits).toHaveLength(2);
    expect(crossReferenceKey(hits[0]!, 0)).not.toBe(crossReferenceKey(hits[1]!, 1));
  });
});
