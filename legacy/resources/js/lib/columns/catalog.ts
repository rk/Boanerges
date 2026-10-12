export type ColumnPickSlot = 'firstNonBible' | 'last';

export type ColumnContentType =
    | 'bible-secondary'
    | 'notes'
    | 'scribe'
    | 'search'
    | 'cross-references'
    | 'comparison'
    | 'verse-list'
    | 'dictionary';

export type ColumnDescriptor = {
    type: ColumnContentType;
    label: string;
    allowDuplicate?: boolean;
    menuId?: string;
    contextMenuLabel?: string;
    pickSlot?: ColumnPickSlot;
};

export const COLUMN_CATALOG: ColumnDescriptor[] = [
    {
        type: 'bible-secondary',
        label: 'Translation',
        allowDuplicate: true,
    },
    { type: 'notes', label: 'Notes' },
    { type: 'scribe', label: 'Scribe' },
    {
        type: 'search',
        label: 'Search',
        menuId: 'study.search',
        pickSlot: 'firstNonBible',
    },
    {
        type: 'cross-references',
        label: 'Cross References',
        menuId: 'study.cross-references',
        contextMenuLabel: 'Cross References',
    },
    {
        type: 'comparison',
        label: 'Comparison',
        menuId: 'study.comparison',
        contextMenuLabel: 'Compare Translations',
    },
    {
        type: 'verse-list',
        label: 'Verse List',
        menuId: 'study.verse-list',
        contextMenuLabel: 'Add to Verse List',
    },
    {
        type: 'dictionary',
        label: 'Dictionary',
        menuId: 'study.dictionary',
        contextMenuLabel: 'Define Word',
    },
];

const catalogByType = new Map<ColumnContentType, ColumnDescriptor>(
    COLUMN_CATALOG.map((descriptor) => [descriptor.type, descriptor]),
);

export function columnDescriptor(type: ColumnContentType): ColumnDescriptor {
    const descriptor = catalogByType.get(type);

    if (!descriptor) {
        throw new Error(`Unknown column type: ${type}`);
    }

    return descriptor;
}

export function columnAllowsDuplicate(type: ColumnContentType): boolean {
    return columnDescriptor(type).allowDuplicate ?? false;
}

export function menuItemForId(
    menuId: string,
): (ColumnDescriptor & { menuId: string }) | undefined {
    const descriptor = COLUMN_CATALOG.find((entry) => entry.menuId === menuId);

    if (descriptor?.menuId === undefined) {
        return undefined;
    }

    return descriptor as ColumnDescriptor & { menuId: string };
}

export type ContextMenuColumnDescriptor = ColumnDescriptor & {
    contextMenuLabel: string;
};

export const CONTEXT_MENU_COLUMNS: ContextMenuColumnDescriptor[] =
    COLUMN_CATALOG.filter(
        (descriptor): descriptor is ContextMenuColumnDescriptor =>
            descriptor.contextMenuLabel !== undefined,
    );
