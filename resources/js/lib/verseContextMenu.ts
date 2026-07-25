import { bible } from '@/lib/bible.svelte.ts';
import { normalizeDictionaryWord } from '@/lib/normalizeDictionaryWord';
import { formatScriptureReference } from '@/lib/scriptureReference';
import {
    ensureComparisonColumn,
    ensureCrossReferencesColumn,
    ensureDictionaryColumn,
    ensureVerseListColumn,
    study,
} from '@/lib/study.svelte.ts';
import { wordAtPoint } from '@/lib/wordAtPoint';

type ContextMenuItem = {
    label: string;
    click?: () => void;
};

export function showVerseContextMenu(
    verse: number,
    event: MouseEvent,
    selectedWord = '',
): void {
    event.preventDefault();

    const reference = formatScriptureReference(
        study.bookId,
        study.chapter,
        verse,
        bible.books,
    );
    const resolvedWord =
        selectedWord.trim() !== '' ? selectedWord : wordAtPoint(event);
    const lookupWord = normalizeDictionaryWord(resolvedWord);
    const items: ContextMenuItem[] = [
        {
            label: 'Cross References',
            click() {
                ensureCrossReferencesColumn(reference);
            },
        },
        {
            label: 'Compare Translations',
            click() {
                ensureComparisonColumn(reference);
            },
        },
        {
            label: 'Add to Verse List',
            click() {
                ensureVerseListColumn({
                    bookId: study.bookId,
                    chapter: study.chapter,
                    verse,
                });
            },
        },
    ];

    if (lookupWord !== '') {
        items.push({
            label: 'Define Word',
            click() {
                ensureDictionaryColumn(lookupWord);
            },
        });
    }

    window.Native?.contextMenu(items satisfies ContextMenuItem[]);
}
