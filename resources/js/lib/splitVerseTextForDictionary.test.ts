import { describe, expect, it } from 'vitest';

import { splitVerseTextForDictionary } from '@/lib/splitVerseTextForDictionary';

describe('splitVerseTextForDictionary', () => {
    it('wraps words separately from punctuation and spaces', () => {
        expect(splitVerseTextForDictionary('the grace, of God')).toEqual([
            { kind: 'word', value: 'the' },
            { kind: 'text', value: ' ' },
            { kind: 'word', value: 'grace' },
            { kind: 'text', value: ', ' },
            { kind: 'word', value: 'of' },
            { kind: 'text', value: ' ' },
            { kind: 'word', value: 'God' },
        ]);
    });

    it('preserves apostrophes inside words', () => {
        expect(splitVerseTextForDictionary("Aaron's rod")).toEqual([
            { kind: 'word', value: "Aaron's" },
            { kind: 'text', value: ' ' },
            { kind: 'word', value: 'rod' },
        ]);
    });
});
