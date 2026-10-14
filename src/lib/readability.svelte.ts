import { Service, type ReadabilitySettings } from "../mygo";
import { readerFontStack } from "./readerFonts";
import { applyReaderTheme } from "./readerTheme";

export const readability = $state({
  fontSize: 16,
  lineHeight: 1.2,
  theme: "auto",
  fontFamily: "serif",
  justifyText: false,
});

let themeCleanup: (() => void) | null = null;

export function syncReaderTheme(): void {
  themeCleanup?.();
  themeCleanup = applyReaderTheme(readability.theme);
}

let hydrated = false;

export function hydrateReadability(settings: ReadabilitySettings): void {
  if (hydrated) return;
  readability.fontSize = settings.fontSize;
  readability.lineHeight = settings.lineHeight;
  readability.theme = settings.theme;
  readability.fontFamily = settings.fontFamily;
  readability.justifyText = settings.justifyText;
  syncReaderTheme();
  hydrated = true;
}

export function readerStyle(): string {
  const family = readerFontStack(readability.fontFamily);
  return `font-size:${readability.fontSize}px;line-height:${readability.lineHeight};font-family:${family};text-align:${readability.justifyText ? "justify" : "start"};`;
}

export async function saveReadability(): Promise<void> {
  const saved = await Service.updateReadability({ ...readability });
  readability.fontSize = saved.fontSize;
  readability.lineHeight = saved.lineHeight;
  readability.theme = saved.theme;
  readability.fontFamily = saved.fontFamily;
  readability.justifyText = saved.justifyText;
  syncReaderTheme();
}
