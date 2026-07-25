import { describe, expect, it } from 'vitest';
import { COLUMN_CATALOG } from '@/lib/columns/catalog';

describe('column menu catalog', () => {
    it('lists Study menu columns in native menu order', () => {
        const menuTypes = COLUMN_CATALOG.filter(
            (descriptor) => descriptor.menuId !== undefined,
        ).map((descriptor) => descriptor.type);

        expect(menuTypes).toEqual([
            'search',
            'cross-references',
            'comparison',
            'verse-list',
            'dictionary',
        ]);
    });
});
