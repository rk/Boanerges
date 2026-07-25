export function verseListEntrySignature(
    entries: ReadonlyArray<{ bookId: string; chapter: number; verse: number }>,
): string {
    return entries
        .map((entry) => `${entry.bookId}:${entry.chapter}:${entry.verse}`)
        .join('|');
}
