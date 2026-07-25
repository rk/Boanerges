import { describe, expect, it } from 'vitest';
import {
    availableColumnOptions,
    columnTargetSlot,
    normalizeColumns,
    sanitizeStudySettings,
} from '@/lib/studyLayout';
import type { ColumnContentType } from '@/lib/types/study';

describe('normalizeColumns', () => {
    it('pads missing slots for two-column layout', () => {
        expect(normalizeColumns(2, [])).toEqual(['bible-secondary']);
    });

    it('trims extra slots when reducing column count', () => {
        const columns: ColumnContentType[] = ['scribe', 'notes'];

        expect(normalizeColumns(2, columns)).toEqual(['scribe']);
    });
});

describe('sanitizeStudySettings', () => {
    it('repairs invalid column count and unknown column types', () => {
        const settings = sanitizeStudySettings({
            columnCount: 4 as 1 | 2 | 3,
            columns: ['invalid', 'notes'] as never,
            bookId: 'gen',
            chapter: 1,
            translationId: 'asv',
            translationBId: 'asv',
            translationCId: 'asv',
        });

        expect(settings.columnCount).toBe(1);
        expect(settings.columns).toEqual([]);
    });

    it('normalizes column slot count for two-column layout', () => {
        const settings = sanitizeStudySettings({
            columnCount: 2,
            columns: ['notes', 'scribe', 'search'],
            bookId: 'gen',
            chapter: 1,
            translationId: 'asv',
            translationBId: 'asv',
            translationCId: 'asv',
        });

        expect(settings.columns).toEqual(['notes']);
    });
});

describe('availableColumnOptions', () => {
    it('prevents duplicate non-bible column types', () => {
        const columns: ColumnContentType[] = ['notes', 'scribe'];
        const options = availableColumnOptions(1, columns, 'asv', 'web');

        expect(options).not.toContain('notes');
        expect(options).toContain('search');
        expect(options).toContain('cross-references');
        expect(options).toContain('dictionary');
        expect(options).toContain('comparison');
        expect(options).toContain('verse-list');
    });

    it('allows two bible columns when translations differ', () => {
        const columns: ColumnContentType[] = ['bible-secondary', 'notes'];
        const options = availableColumnOptions(1, columns, 'asv', 'web');

        expect(options).toContain('bible-secondary');
    });
});

describe('columnTargetSlot', () => {
    it('prefers the last slot already showing the requested type', () => {
        expect(
            columnTargetSlot('cross-references', 3, [
                'search',
                'cross-references',
            ]),
        ).toBe(1);
        expect(
            columnTargetSlot('cross-references', 3, [
                'cross-references',
                'notes',
            ]),
        ).toBe(0);
        expect(
            columnTargetSlot('dictionary', 3, ['search', 'dictionary']),
        ).toBe(1);
        expect(columnTargetSlot('comparison', 3, ['notes', 'comparison'])).toBe(
            1,
        );
        expect(
            columnTargetSlot('verse-list', 3, ['search', 'verse-list']),
        ).toBe(1);
    });

    it('falls back to the last secondary slot', () => {
        expect(columnTargetSlot('cross-references', 2, ['notes'])).toBe(0);
        expect(columnTargetSlot('dictionary', 3, ['notes', 'search'])).toBe(1);
    });
});
