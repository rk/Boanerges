import type { ColumnContentType, StudySettings } from '@/lib/types/study';

export const COLUMN_CONTENT_TYPES: ColumnContentType[] = [
    'bible-secondary',
    'notes',
    'scribe',
    'search',
    'cross-references',
];

export const COLUMN_CONTENT_LABELS: Record<ColumnContentType, string> = {
    'bible-secondary': 'Translation',
    notes: 'Notes',
    scribe: 'Scribe',
    search: 'Search',
    'cross-references': 'Cross References',
};

export function isColumnCount(value: number): value is 1 | 2 | 3 {
    return value === 1 || value === 2 || value === 3;
}

export function sanitizeColumnType(value: unknown): ColumnContentType | null {
    if (typeof value !== 'string') {
        return null;
    }

    return COLUMN_CONTENT_TYPES.includes(value as ColumnContentType)
        ? (value as ColumnContentType)
        : null;
}

export function sanitizeStudySettings(settings: StudySettings): StudySettings {
    const columnCount = isColumnCount(settings.columnCount)
        ? settings.columnCount
        : 1;
    const columns = (settings.columns ?? [])
        .map((column) => sanitizeColumnType(column))
        .filter((column): column is ColumnContentType => column !== null);

    return {
        ...settings,
        columnCount,
        columns: normalizeColumns(columnCount, columns),
    };
}

export function normalizeColumns(
    columnCount: 1 | 2 | 3,
    columns: ColumnContentType[],
): ColumnContentType[] {
    const needed = Math.max(0, columnCount - 1);
    const normalized = columns.slice(0, needed);

    while (normalized.length < needed) {
        normalized.push(defaultColumnForSlot(normalized));
    }

    return normalized;
}

function defaultColumnForSlot(
    existing: ColumnContentType[],
): ColumnContentType {
    for (const type of COLUMN_CONTENT_TYPES) {
        if (type === 'bible-secondary' || !existing.includes(type)) {
            return type;
        }
    }

    return 'bible-secondary';
}

export function availableColumnOptions(
    slotIndex: number,
    columns: ColumnContentType[],
    translationBId: string,
    translationCId: string,
): ColumnContentType[] {
    const used = new Set(columns.filter((_, index) => index !== slotIndex));

    return COLUMN_CONTENT_TYPES.filter((type) => {
        if (type !== 'bible-secondary' && used.has(type)) {
            return false;
        }

        if (type !== 'bible-secondary') {
            return true;
        }

        const otherBibleSlots = columns
            .map((column, index) => ({ column, index }))
            .filter(
                ({ column, index }) =>
                    column === 'bible-secondary' && index !== slotIndex,
            );

        if (otherBibleSlots.length === 0) {
            return true;
        }

        const otherIndex = otherBibleSlots[0].index;
        const otherTranslation =
            otherIndex === 0 ? translationBId : translationCId;
        const thisTranslation =
            slotIndex === 0 ? translationBId : translationCId;

        return otherTranslation !== thisTranslation;
    });
}

export function bibleColumnCount(
    settings: Pick<StudySettings, 'columnCount' | 'columns'>,
): number {
    return (
        1 +
        settings.columns.filter((column) => column === 'bible-secondary').length
    );
}

export function secondaryTranslationForSlot(
    slotIndex: number,
    translationBId: string,
    translationCId: string,
): string {
    return slotIndex === 0 ? translationBId : translationCId;
}

export function setColumnContentType(
    settings: StudySettings,
    slotIndex: number,
    type: ColumnContentType,
): StudySettings {
    const columns = [...settings.columns];
    columns[slotIndex] = type;

    return {
        ...settings,
        columns: normalizeColumns(settings.columnCount, columns),
    };
}

export function crossReferencesTargetSlot(
    columnCount: 1 | 2 | 3,
    columns: ColumnContentType[],
): number {
    for (let index = columns.length - 1; index >= 0; index--) {
        if (columns[index] === 'cross-references') {
            return index;
        }
    }

    return Math.max(0, columnCount - 2);
}
