import { bible } from '@/lib/bible.svelte.ts';
import { formatScriptureReference } from '@/lib/scriptureReference';
import { ensureCrossReferencesColumn, study } from '@/lib/study.svelte.ts';

type ContextMenuItem = {
    label: string;
    click?: () => void;
};

export function showVerseContextMenu(verse: number, event: MouseEvent): void {
    event.preventDefault();

    const reference = formatScriptureReference(
        study.bookId,
        study.chapter,
        verse,
        bible.books,
    );

    window.Native?.contextMenu([
        {
            label: 'Cross References',
            click() {
                ensureCrossReferencesColumn(reference);
            },
        },
    ] satisfies ContextMenuItem[]);
}
