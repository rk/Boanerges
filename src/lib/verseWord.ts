/** First dictionary word from selection or verse plain text. */
export function wordFromVerseContext(event: MouseEvent, verseTextHtml: string): string {
  const selected = window.getSelection()?.toString().trim();
  if (selected && /[A-Za-z]/.test(selected)) {
    const m = selected.match(/[A-Za-z][A-Za-z'-]*/);
    if (m) {
      return m[0].toLowerCase();
    }
  }
  const plain = verseTextHtml.replace(/<[^>]+>/g, " ");
  const m = plain.match(/[A-Za-z][A-Za-z'-]{1,}/);
  return m ? m[0].toLowerCase() : "";
}

export function clampMenuPosition(x: number, y: number, menuWidth = 208, menuHeight = 200): { x: number; y: number } {
  const pad = 8;
  const maxX = window.innerWidth - menuWidth - pad;
  const maxY = window.innerHeight - menuHeight - pad;
  return {
    x: Math.max(pad, Math.min(x, maxX)),
    y: Math.max(pad, Math.min(y, maxY)),
  };
}
