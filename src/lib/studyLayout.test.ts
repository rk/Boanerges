import { describe, expect, it } from "vitest";
import { normalizeColumns } from "./studyLayout";

describe("normalizeColumns", () => {
  it("pads secondary slots when column count increases", () => {
    expect(normalizeColumns(2, [])).toEqual(["bible-secondary"]);
    expect(normalizeColumns(3, [])).toEqual(["bible-secondary", "bible-secondary"]);
  });
});
