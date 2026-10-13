import { describe, expect, it } from "vitest";
import { normalizeDictionaryWord } from "./normalizeDictionaryWord";

describe("normalizeDictionaryWord", () => {
  it("takes first token and strips punctuation", () => {
    expect(normalizeDictionaryWord("  earth, ")).toBe("earth");
    expect(normalizeDictionaryWord("God said")).toBe("God");
  });
});
