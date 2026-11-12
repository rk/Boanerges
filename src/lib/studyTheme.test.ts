import { describe, expect, it } from "vitest";
import { studyThemeName } from "./studyTheme";

describe("studyThemeName", () => {
  it("composes pole and appearance", () => {
    expect(studyThemeName("modern", "light")).toBe("study-modern-light");
    expect(studyThemeName("monastic", "dark")).toBe("study-monastic-dark");
  });
});
