import type { Component } from 'svelte';
import ComparisonColumn from '@/components/columns/ComparisonColumn.svelte';
import CrossReferencesColumn from '@/components/columns/CrossReferencesColumn.svelte';
import DictionaryColumn from '@/components/columns/DictionaryColumn.svelte';
import NotesColumn from '@/components/columns/NotesColumn.svelte';
import ScribeColumn from '@/components/columns/ScribeColumn.svelte';
import SearchColumn from '@/components/columns/SearchColumn.svelte';
import VerseListColumn from '@/components/columns/VerseListColumn.svelte';
import type { ColumnContentType } from '@/lib/columns/catalog';

type ColumnComponent = Component<{ slotIndex: number }>;

export const COLUMN_COMPONENTS: Record<
    Exclude<ColumnContentType, 'bible-secondary'>,
    ColumnComponent
> = {
    notes: NotesColumn,
    scribe: ScribeColumn,
    search: SearchColumn,
    'cross-references': CrossReferencesColumn,
    dictionary: DictionaryColumn,
    comparison: ComparisonColumn,
    'verse-list': VerseListColumn,
};
