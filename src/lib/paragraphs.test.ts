import { describe, expect, it } from "vitest";
import { groupVersesIntoParagraphs } from "./paragraphs";
import type { Verse } from "../mygo";

describe("groupVersesIntoParagraphs", () => {
  it("keeps consecutive verses in one group until paragraphStart", () => {
    const verses: Verse[] = [
      { number: 1, text: "A", paragraphStart: true },
      { number: 2, text: "B" },
      { number: 3, text: "C", paragraphStart: true },
      { number: 4, text: "D" },
    ];
    expect(groupVersesIntoParagraphs(verses)).toEqual([
      [
        { number: 1, text: "A", paragraphStart: true },
        { number: 2, text: "B" },
      ],
      [
        { number: 3, text: "C", paragraphStart: true },
        { number: 4, text: "D" },
      ],
    ]);
  });
});
