import { describe, expect, it } from 'vitest';
import { normalizeDictionaryWord } from '@/lib/normalizeDictionaryWord';

describe('normalizeDictionaryWord', () => {
    it('returns the first token from a selection', () => {
        expect(normalizeDictionaryWord('grace and mercy')).toBe('grace');
    });

    it('strips leading and trailing punctuation', () => {
        expect(normalizeDictionaryWord('"salvation,"')).toBe('salvation');
    });

    it('preserves hyphens and apostrophes', () => {
        expect(normalizeDictionaryWord("'AARON'S'")).toBe("'AARON'S'");
    });

    it('returns empty string for whitespace-only input', () => {
        expect(normalizeDictionaryWord('   ')).toBe('');
    });
});
