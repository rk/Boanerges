import { describe, expect, it } from 'vitest';

import { verseNumbersInRange } from '@/lib/verseHighlight';

describe('verseNumbersInRange', () => {
    it('returns a single verse when end is omitted', () => {
        expect([...verseNumbersInRange(3, null)]).toEqual([3]);
    });

    it('returns an inclusive range when end is provided', () => {
        expect([...verseNumbersInRange(3, 7)]).toEqual([3, 4, 5, 6, 7]);
    });
});
