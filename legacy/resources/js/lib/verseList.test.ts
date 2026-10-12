import { describe, expect, it } from 'vitest';
import { verseListEntrySignature } from '@/lib/verseListSignature';
import { verseListNeedsTextLoad } from '@/lib/verseListTextLoad';

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
    const exists = entries.some((entry) => entryKey(entry) === entryKey(ref));

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

describe('verseListEntrySignature', () => {
    it('ignores verse text when building signature', () => {
        const entries = [
            { bookId: 'gen', chapter: 1, verse: 1, text: 'Before' },
            { bookId: 'jhn', chapter: 3, verse: 16, text: 'After' },
        ];

        expect(verseListEntrySignature(entries)).toBe('gen:1:1|jhn:3:16');
    });
});

describe('verseListNeedsTextLoad', () => {
    it('loads when a verse has no text', () => {
        expect(verseListNeedsTextLoad([{}], true, 'asv', false, null)).toBe(
            true,
        );
    });

    it('skips while loading', () => {
        expect(verseListNeedsTextLoad([{}], true, 'asv', true, null)).toBe(
            false,
        );
    });

    it('reloads when translation changes', () => {
        expect(
            verseListNeedsTextLoad(
                [{ text: 'In the beginning' }],
                true,
                'kjv',
                false,
                'asv',
            ),
        ).toBe(true);

        expect(
            verseListNeedsTextLoad(
                [{ text: 'In the beginning' }],
                true,
                'asv',
                false,
                'asv',
            ),
        ).toBe(false);
    });
});
