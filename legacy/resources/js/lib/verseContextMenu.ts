import { bible } from '@/lib/bible.svelte.ts';
import { CONTEXT_MENU_COLUMNS } from '@/lib/columns/catalog';
import { normalizeDictionaryWord } from '@/lib/normalizeDictionaryWord';
import { formatScriptureReference } from '@/lib/scriptureReference';
import { openColumn, study } from '@/lib/study.svelte.ts';
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
    const items: ContextMenuItem[] = CONTEXT_MENU_COLUMNS.filter(
        (descriptor) => descriptor.type !== 'dictionary' || lookupWord !== '',
    ).map((descriptor) => ({
        label: descriptor.contextMenuLabel,
        click() {
            if (descriptor.type === 'cross-references') {
                openColumn('cross-references', {
                    kind: 'reference',
                    reference,
                });
            } else if (descriptor.type === 'comparison') {
                openColumn('comparison', {
                    kind: 'reference',
                    reference,
                });
            } else if (descriptor.type === 'verse-list') {
                openColumn('verse-list', {
                    kind: 'verse',
                    ref: {
                        bookId: study.bookId,
                        chapter: study.chapter,
                        verse,
                    },
                });
            } else if (descriptor.type === 'dictionary') {
                openColumn('dictionary', {
                    kind: 'word',
                    word: lookupWord,
                });
            }
        },
    }));

    window.Native?.contextMenu(items satisfies ContextMenuItem[]);
}
