export type ReaderFontFamily = "sans-serif" | "serif";

/** Match legacy reader stacks (see style.css --font-*). */
export const readerFontStacks: Record<ReaderFontFamily, string> = {
  "sans-serif": 'system-ui, -apple-system, "Segoe UI", sans-serif',
  serif: '"Iowan Old Style", "Palatino Linotype", Palatino, Georgia, serif',
};

export function readerFontStack(family: string): string {
  if (family === "sans-serif") {
    return readerFontStacks["sans-serif"];
  }
  return readerFontStacks.serif;
}

export function isReaderFontFamily(value: string): value is ReaderFontFamily {
  return value === "sans-serif" || value === "serif";
}
