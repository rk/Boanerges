import { describe, expect, it } from 'vitest';
import {
    activeBibleTranslationIds,
    isBookAvailableInTranslations,
} from '@/lib/studyLayout';
import type { Book } from '@/lib/types/bible';

const fullBible: Book[] = [
    {
        id: 'gen',
        name: 'Genesis',
        abbrev: 'GEN',
        testament: 'ot',
        chapters: 50,
    },
    {
        id: 'mat',
        name: 'Matthew',
        abbrev: 'MAT',
        testament: 'nt',
        chapters: 28,
    },
];

const ntOnly: Book[] = [
    {
        id: 'mat',
        name: 'Matthew',
        abbrev: 'MAT',
        testament: 'nt',
        chapters: 28,
    },
];

describe('activeBibleTranslationIds', () => {
    it('returns only primary translation for single-column layout', () => {
        expect(
            activeBibleTranslationIds({
                translationId: 'asv',
                translationBId: 'web',
                translationCId: 'kjv',
                columns: [],
            }),
        ).toEqual(['asv']);
    });

    it('includes secondary bible column translations', () => {
        expect(
            activeBibleTranslationIds({
                translationId: 'asv',
                translationBId: 'web',
                translationCId: 'kjv',
                columns: ['bible-secondary', 'notes'],
            }),
        ).toEqual(['asv', 'web']);
    });

    it('includes all bible columns in three-column layout', () => {
        expect(
            activeBibleTranslationIds({
                translationId: 'asv',
                translationBId: 'web',
                translationCId: 'kjv',
                columns: ['bible-secondary', 'bible-secondary'],
            }),
        ).toEqual(['asv', 'web', 'kjv']);
    });

    it('deduplicates when secondary slots share the same translation', () => {
        expect(
            activeBibleTranslationIds({
                translationId: 'asv',
                translationBId: 'web',
                translationCId: 'web',
                columns: ['bible-secondary', 'bible-secondary'],
            }),
        ).toEqual(['asv', 'web']);
    });

    it('ignores non-bible secondary columns', () => {
        expect(
            activeBibleTranslationIds({
                translationId: 'asv',
                translationBId: 'web',
                translationCId: 'kjv',
                columns: ['notes', 'scribe'],
            }),
        ).toEqual(['asv']);
    });
});

describe('isBookAvailableInTranslations', () => {
    const booksByTranslation = new Map<string, readonly Book[]>([
        ['asv', fullBible],
        ['web', fullBible],
        ['nt', ntOnly],
    ]);

    it('returns true when every translation has the book', () => {
        expect(
            isBookAvailableInTranslations(
                'gen',
                ['asv', 'web'],
                booksByTranslation,
            ),
        ).toBe(true);
    });

    it('returns false when any translation lacks the book', () => {
        expect(
            isBookAvailableInTranslations(
                'gen',
                ['asv', 'nt'],
                booksByTranslation,
            ),
        ).toBe(false);
    });

    it('returns false when no translation has the book', () => {
        expect(
            isBookAvailableInTranslations('rev', ['nt'], booksByTranslation),
        ).toBe(false);
    });

    it('returns true for a book present in a single translation', () => {
        expect(
            isBookAvailableInTranslations('mat', ['nt'], booksByTranslation),
        ).toBe(true);
    });
});
