import { bible, fetchChapter } from '@/lib/bible.svelte.ts';
import { parseScriptureReference } from '@/lib/scriptureReference';
import type { ScriptureReference } from '@/lib/scriptureReference';
import type { Book } from '@/lib/types/bible';

export type ComparisonRow = {
    translationId: string;
    name: string;
    abbrev: string;
    text: string;
};

export const comparison = $state({
    loading: false,
    activeReference: null as string | null,
    rows: [] as ComparisonRow[],
});

let debounceTimeout: ReturnType<typeof setTimeout> | null = null;
let activeRequest = 0;

export function setComparisonInput(reference: string): void {
    comparison.activeReference = reference;
}

export function scheduleComparisonLookup(input: string, books: Book[]): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
    }

    const trimmed = input.trim();

    if (trimmed === '') {
        activeRequest++;
        comparison.loading = false;
        comparison.rows = [];

        return;
    }

    const parsed = parseScriptureReference(trimmed, books);

    if (parsed === null) {
        activeRequest++;
        comparison.loading = false;
        comparison.rows = [];

        return;
    }

    debounceTimeout = setTimeout(() => {
        void loadComparison(parsed);
    }, 300);
}

export async function loadComparison(ref: ScriptureReference): Promise<void> {
    const requestId = ++activeRequest;
    comparison.loading = true;

    try {
        const rows = await Promise.all(
            bible.translations.map(async (translation) => {
                try {
                    const chapter = await fetchChapter(
                        translation.id,
                        ref.bookId,
                        ref.chapter,
                    );
                    const verse = chapter.verses.find(
                        (entry) => entry.number === ref.verse,
                    );

                    return {
                        translationId: translation.id,
                        name: translation.name,
                        abbrev: translation.abbrev,
                        text: verse?.text ?? '—',
                    };
                } catch {
                    return {
                        translationId: translation.id,
                        name: translation.name,
                        abbrev: translation.abbrev,
                        text: '—',
                    };
                }
            }),
        );

        if (requestId !== activeRequest) {
            return;
        }

        comparison.rows = rows;
    } finally {
        if (requestId === activeRequest) {
            comparison.loading = false;
        }
    }
}
