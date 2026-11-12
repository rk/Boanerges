import { Service, type ReadabilitySettings } from "../mygo";
import { readerFontStack } from "./readerFonts";
import { applyReadabilityTheme } from "./readerTheme";
import type { StudyAppearance, StudyThemePole } from "./studyTheme";

export const readability = $state({
  fontSize: 17,
  lineHeight: 1.6,
  theme: "auto" as StudyAppearance,
  themePole: "monastic" as StudyThemePole,
  fontFamily: "serif" as "serif" | "sans-serif",
  justifyText: true,
});

let themeCleanup: (() => void) | null = null;

export function syncReaderTheme(): void {
  themeCleanup?.();
  themeCleanup = applyReadabilityTheme(readability.themePole, readability.theme);
}

let hydrated = false;

export function hydrateReadability(settings: ReadabilitySettings): void {
  if (hydrated) return;
  readability.fontSize = settings.fontSize;
  readability.lineHeight = settings.lineHeight;
  const theme = settings.theme;
  readability.theme =
    theme === "light" || theme === "dark" || theme === "auto" ? theme : "auto";
  const pole = settings.themePole;
  readability.themePole = pole === "modern" ? "modern" : "monastic";
  readability.fontFamily = settings.fontFamily === "sans-serif" ? "sans-serif" : "serif";
  readability.justifyText = settings.justifyText;
  syncReaderTheme();
  hydrated = true;
}

export function readerStyle(): string {
  const family = readerFontStack(readability.fontFamily);
  return `font-size:${readability.fontSize}px;line-height:${readability.lineHeight};font-family:${family};text-align:${readability.justifyText ? "justify" : "start"};`;
}

export async function saveReadability(): Promise<void> {
  const saved = await Service.updateReadability({
    fontSize: readability.fontSize,
    lineHeight: readability.lineHeight,
    theme: readability.theme,
    themePole: readability.themePole,
    fontFamily: readability.fontFamily,
    justifyText: readability.justifyText,
  });
  readability.fontSize = saved.fontSize;
  readability.lineHeight = saved.lineHeight;
  readability.theme =
    saved.theme === "light" || saved.theme === "dark" || saved.theme === "auto"
      ? saved.theme
      : "auto";
  readability.themePole = saved.themePole === "modern" ? "modern" : "monastic";
  readability.fontFamily = saved.fontFamily === "sans-serif" ? "sans-serif" : "serif";
  readability.justifyText = saved.justifyText;
  syncReaderTheme();
}
