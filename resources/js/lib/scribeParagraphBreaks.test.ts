import { describe, expect, it } from 'vitest';

import {
    applyPrimaryParagraphBreaks,
    effectiveParagraphStart,
    resetInterVerseParagraphBreaks,
} from '@/lib/scribeParagraphBreaks.ts';
import type { Verse } from '@/lib/types/bible';

describe('effectiveParagraphStart', () => {
    it('returns override when set', () => {
        expect(effectiveParagraphStart(2, true)).toBe(true);
        expect(effectiveParagraphStart(2, false)).toBe(false);
    });

    it('defaults verse 1 to true and later verses to false', () => {
        expect(effectiveParagraphStart(1)).toBe(true);
        expect(effectiveParagraphStart(2)).toBe(false);
        expect(effectiveParagraphStart(10)).toBe(false);
    });
});

describe('applyPrimaryParagraphBreaks', () => {
    const sourceVerses: Verse[] = [
        { number: 1, text: 'A', paragraphStart: true },
        { number: 2, text: 'B' },
        { number: 3, text: 'C', paragraphStart: true },
    ];

    it('copies primary paragraph flags without changing verse text', () => {
        const entries = {
            2: { text: 'User verse two\n\nwith breaks' },
            3: { text: 'User verse three', paragraphStartOverride: false },
        };

        const result = applyPrimaryParagraphBreaks(
            [1, 2, 3],
            entries,
            sourceVerses,
        );

        expect(result[2]).toEqual({
            text: 'User verse two\n\nwith breaks',
        });
        expect(result[3]).toEqual({
            text: 'User verse three',
            paragraphStartOverride: true,
        });
    });

    it('overwrites prior inter-verse overrides from primary', () => {
        const entries = {
            2: { text: 'Two', paragraphStartOverride: true },
            3: { text: 'Three', paragraphStartOverride: true },
        };

        const result = applyPrimaryParagraphBreaks(
            [1, 2, 3],
            entries,
            sourceVerses,
        );

        expect(result[2]).toEqual({ text: 'Two' });
        expect(result[3]).toEqual({
            text: 'Three',
            paragraphStartOverride: true,
        });
    });
});

describe('resetInterVerseParagraphBreaks', () => {
    it('strips inter-verse overrides but preserves intra-verse text', () => {
        const entries = {
            2: { text: 'Line one\n\nLine two', paragraphStartOverride: true },
            5: { text: '', paragraphStartOverride: true },
            7: { text: 'Only text', paragraphStartOverride: false },
        };

        const result = resetInterVerseParagraphBreaks([1, 2, 3, 5, 7], entries);

        expect(result[2]).toEqual({ text: 'Line one\n\nLine two' });
        expect(result[5]).toBeUndefined();
        expect(result[7]).toEqual({ text: 'Only text' });
    });
});
