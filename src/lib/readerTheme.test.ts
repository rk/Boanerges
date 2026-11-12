import { describe, expect, it } from "vitest";
import { applyReadabilityTheme } from "./readerTheme";

describe("applyReadabilityTheme", () => {
  it("sets composed study theme on documentElement", () => {
    applyReadabilityTheme("monastic", "light")();
    expect(document.documentElement.getAttribute("data-theme")).toBe("study-monastic-light");
    applyReadabilityTheme("modern", "dark")();
    expect(document.documentElement.getAttribute("data-theme")).toBe("study-modern-dark");
  });

  it("auto uses light or dark from matchMedia", () => {
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    applyReadabilityTheme("monastic", "auto")();
    const expected = mq.matches ? "study-monastic-dark" : "study-monastic-light";
    expect(document.documentElement.getAttribute("data-theme")).toBe(expected);
  });
});
