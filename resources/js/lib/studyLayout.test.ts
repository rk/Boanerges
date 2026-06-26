import { describe, expect, it } from 'vitest';
import {
    availableColumnOptions,
    crossReferencesTargetSlot,
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
    });

    it('allows two bible columns when translations differ', () => {
        const columns: ColumnContentType[] = ['bible-secondary', 'notes'];
        const options = availableColumnOptions(1, columns, 'asv', 'web');

        expect(options).toContain('bible-secondary');
    });
});

describe('crossReferencesTargetSlot', () => {
    it('prefers the last slot already showing cross references', () => {
        expect(crossReferencesTargetSlot(3, ['search', 'cross-references'])).toBe(1);
        expect(crossReferencesTargetSlot(3, ['cross-references', 'notes'])).toBe(0);
    });

    it('falls back to the last secondary slot', () => {
        expect(crossReferencesTargetSlot(2, ['notes'])).toBe(0);
        expect(crossReferencesTargetSlot(3, ['notes', 'search'])).toBe(1);
    });
});
