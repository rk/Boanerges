import { Service, type ReadabilitySettings } from "../mygo";

export const readability = $state({
  fontSize: 18,
  lineHeight: 1.6,
  theme: "auto",
  fontFamily: "serif",
  justifyText: false,
});

let hydrated = false;

export function hydrateReadability(settings: ReadabilitySettings): void {
  if (hydrated) return;
  readability.fontSize = settings.fontSize;
  readability.lineHeight = settings.lineHeight;
  readability.theme = settings.theme;
  readability.fontFamily = settings.fontFamily;
  readability.justifyText = settings.justifyText;
  hydrated = true;
}

export function readerStyle(): string {
  const family = readability.fontFamily === "serif" ? "Georgia, serif" : "system-ui, sans-serif";
  return `font-size:${readability.fontSize}px;line-height:${readability.lineHeight};font-family:${family};text-align:${readability.justifyText ? "justify" : "start"};`;
}

export async function saveReadability(): Promise<void> {
  await Service.updateReadability({ ...readability });
}
