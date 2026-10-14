import { describe, expect, it } from "vitest";
import { applyReaderTheme } from "./readerTheme";

describe("applyReaderTheme", () => {
  it("sets light and dark on html", () => {
    applyReaderTheme("light")();
    expect(document.documentElement.getAttribute("data-theme")).toBe("light");
    applyReaderTheme("dark")();
    expect(document.documentElement.getAttribute("data-theme")).toBe("dark");
  });

  it("removes data-theme for auto", () => {
    document.documentElement.setAttribute("data-theme", "dark");
    applyReaderTheme("auto")();
    expect(document.documentElement.hasAttribute("data-theme")).toBe(false);
  });
});
