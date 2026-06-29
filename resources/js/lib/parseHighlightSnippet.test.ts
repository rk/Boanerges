import { describe, expect, it } from 'vitest';

import { parseHighlightSnippet } from '@/lib/parseHighlightSnippet';

describe('parseHighlightSnippet', () => {
    it('returns empty array for empty snippet', () => {
        expect(parseHighlightSnippet('')).toEqual([]);
    });

    it('returns plain text when no mark tags', () => {
        expect(parseHighlightSnippet('In the beginning')).toEqual([
            { text: 'In the beginning', highlight: false },
        ]);
    });

    it('splits on mark tags', () => {
        expect(
            parseHighlightSnippet('… the <mark>word</mark> of God …'),
        ).toEqual([
            { text: '… the ', highlight: false },
            { text: 'word', highlight: true },
            { text: ' of God …', highlight: false },
        ]);
    });

    it('handles multiple highlights', () => {
        expect(
            parseHighlightSnippet('<mark>God</mark> and <mark>light</mark>'),
        ).toEqual([
            { text: 'God', highlight: true },
            { text: ' and ', highlight: false },
            { text: 'light', highlight: true },
        ]);
    });
});
