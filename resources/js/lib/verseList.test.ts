import { describe, expect, it } from 'vitest';

function entryKey(entry: {
    bookId: string;
    chapter: number;
    verse: number;
}): string {
    return `${entry.bookId}:${entry.chapter}:${entry.verse}`;
}

function dedupeEntries(
    entries: Array<{ bookId: string; chapter: number; verse: number }>,
    ref: { bookId: string; chapter: number; verse: number },
): Array<{ bookId: string; chapter: number; verse: number }> {
    const exists = entries.some(
        (entry) => entryKey(entry) === entryKey(ref),
    );

    if (exists) {
        return entries;
    }

    return [...entries, ref];
}

describe('verse list dedupe', () => {
    it('skips duplicate verses', () => {
        const entries = [{ bookId: 'gen', chapter: 1, verse: 1 }];

        expect(
            dedupeEntries(entries, { bookId: 'gen', chapter: 1, verse: 1 }),
        ).toEqual(entries);
    });

    it('appends new verses', () => {
        const entries = [{ bookId: 'gen', chapter: 1, verse: 1 }];

        expect(
            dedupeEntries(entries, { bookId: 'jhn', chapter: 3, verse: 16 }),
        ).toEqual([
            { bookId: 'gen', chapter: 1, verse: 1 },
            { bookId: 'jhn', chapter: 3, verse: 16 },
        ]);
    });
});
