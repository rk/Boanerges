import { bible } from '@/lib/bible.svelte.ts';
import { normalizeDictionaryWord } from '@/lib/normalizeDictionaryWord';
import { formatScriptureReference } from '@/lib/scriptureReference';
import {
    ensureCrossReferencesColumn,
    ensureDictionaryColumn,
    study,
} from '@/lib/study.svelte.ts';

type ContextMenuItem = {
    label: string;
    click?: () => void;
};

function wordAtPoint(event: MouseEvent): string {
    const range = document.caretRangeFromPoint(event.clientX, event.clientY);

    if (!range || !range.startContainer.textContent) {
        return '';
    }

    const text = range.startContainer.textContent;
    const offset = range.startOffset;
    const before = text.slice(0, offset);
    const after = text.slice(offset);
    const start = before.search(/[\p{L}\p{N}'-]+$/u);
    const endMatch = after.match(/^[\p{L}\p{N}'-]+/u);
    const end = endMatch ? offset + endMatch[0].length : offset;

    if (start === -1) {
        return '';
    }

    return text.slice(start, end);
}

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
