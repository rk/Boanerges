import { describe, expect, it } from "vitest";
import {
  focusAdjacentVerse,
  hydrateDocument,
  hydrateVerseSpan,
  normalizePasteText,
  paragraphStartAttribute,
  readVerseSpan,
  serializeFromDocument,
  shouldHydrateVerseSpan,
} from "./scribeDocument";

describe("normalizePasteText", () => {
  it("normalizes Windows and classic Mac line endings", () => {
    expect(normalizePasteText("a\r\nb\rc")).toBe("a\nb\nc");
  });
});

describe("verse span hydration", () => {
  it("skips hydration while focused unless forced", () => {
    const span = document.createElement("span");
    span.textContent = "live draft";
    document.body.appendChild(span);
    span.focus();

    expect(shouldHydrateVerseSpan(span, "server draft")).toBe(false);

    hydrateVerseSpan(span, "server draft", { force: true });

    expect(readVerseSpan(span)).toBe("server draft");
  });
});

describe("document serialization", () => {
  it("round-trips verse text", () => {
    const root = document.createElement("div");
    root.innerHTML = `
      <span class="scribe-verse" data-verse="1">First line</span>
      <span class="scribe-verse" data-verse="3">Third verse</span>
    `;

    expect(serializeFromDocument(root)).toEqual({
      1: { text: "First line" },
      3: { text: "Third verse" },
    });
  });
});

describe("paragraphStartAttribute", () => {
  it("maps booleans to data attribute values", () => {
    expect(paragraphStartAttribute(true)).toBe("true");
    expect(paragraphStartAttribute(false)).toBe("false");
  });
});

describe("focusAdjacentVerse", () => {
  it("moves focus between verse spans", () => {
    const root = document.createElement("div");
    root.innerHTML = `
      <span class="scribe-verse" data-verse="1" tabindex="0"></span>
      <span class="scribe-verse" data-verse="2" tabindex="0"></span>
    `;

    document.body.appendChild(root);

    const spans = [...root.querySelectorAll<HTMLElement>(".scribe-verse")];
    spans[0]!.focus();

    expect(focusAdjacentVerse(root, 1, 1)).toBe(true);
    expect(document.activeElement).toBe(spans[1]!);
  });
});

describe("hydrateDocument", () => {
  it("hydrates all verse spans from entries", () => {
    const root = document.createElement("div");
    root.innerHTML = `
      <span class="scribe-verse-wrap"><span class="scribe-verse" data-verse="1"></span></span>
      <span class="scribe-verse-wrap"><span class="scribe-verse" data-verse="2"></span></span>
    `;

    hydrateDocument(
      root,
      { 1: { text: "Alpha" }, 2: { text: "Beta" } },
      { 1: true, 2: false },
      { force: true },
    );

    const spans = [...root.querySelectorAll<HTMLElement>(".scribe-verse")];

    expect(readVerseSpan(spans[0]!)).toBe("Alpha");
    expect(spans[0]!.dataset.paragraphStart).toBe("true");
  });
});
