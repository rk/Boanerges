import { describe, expect, it } from 'vitest';

import { parseScriptureReference, formatScriptureReference } from '@/lib/scriptureReference';
import type { Book } from '@/lib/types/bible';

const books: Book[] = [
    { id: 'gen', name: 'Genesis', abbrev: 'GEN', testament: 'ot', chapters: 50 },
    { id: 'mrk', name: 'Mark', abbrev: 'MRK', testament: 'nt', chapters: 16 },
    { id: '1jn', name: '1 John', abbrev: '1JN', testament: 'nt', chapters: 5 },
    { id: 'rut', name: 'Ruth', abbrev: 'RUT', testament: 'ot', chapters: 4 },
];

describe('parseScriptureReference', () => {
    it('parses a single-word book name', () => {
        expect(parseScriptureReference('Mark 1:1', books)).toEqual({
            bookId: 'mrk',
            chapter: 1,
            verse: 1,
        });
    });

    it('parses numbered book names', () => {
        expect(parseScriptureReference('1 John 3:16', books)).toEqual({
            bookId: '1jn',
            chapter: 3,
            verse: 16,
        });
    });

    it('matches abbreviations case-insensitively', () => {
        expect(parseScriptureReference('rut 1:1', books)).toEqual({
            bookId: 'rut',
            chapter: 1,
            verse: 1,
        });
    });

    it('returns null for invalid references', () => {
        expect(parseScriptureReference('Not A Book 1:1', books)).toBeNull();
        expect(parseScriptureReference('Mark 1', books)).toBeNull();
    });
});

describe('formatScriptureReference', () => {
    it('uses canon names when the translation catalog is unavailable', () => {
        expect(formatScriptureReference('ezr', 1, 1, [])).toBe('Ezra 1:1');
        expect(formatScriptureReference('zec', 4, 6, [])).toBe('Zechariah 4:6');
    });
});
