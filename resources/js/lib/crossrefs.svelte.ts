import { crossReferences as crossReferencesRoute } from '@/actions/App/Http/Controllers/BibleController';
import { formatScriptureReference, parseScriptureReference } from '@/lib/scriptureReference';
import { study } from '@/lib/study.svelte.ts';
import type { Book, CrossReference } from '@/lib/types/bible';

export const crossrefs = $state({
    open: false,
    loading: false,
    references: [] as CrossReference[],
    initialReference: null as string | null,
});

let debounceTimeout: ReturnType<typeof setTimeout> | null = null;
let activeRequest = 0;

export function scriptureReferenceForCurrentVerse(books: Book[]): string {
    const verse = study.verseHighlight?.verse ?? 1;

    return formatScriptureReference(study.bookId, study.chapter, verse, books);
}

export function openCrossReferences(books: Book[], reference?: string): void {
    crossrefs.initialReference = reference ?? scriptureReferenceForCurrentVerse(books);
    crossrefs.open = true;
}

export function closeCrossReferences(): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
        debounceTimeout = null;
    }

    activeRequest++;
    crossrefs.open = false;
    crossrefs.loading = false;
    crossrefs.references = [];
    crossrefs.initialReference = null;
}

export function scheduleCrossReferenceLookup(input: string, books: Book[]): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
    }

    const trimmed = input.trim();

    if (trimmed === '') {
        activeRequest++;
        crossrefs.loading = false;
        crossrefs.references = [];

        return;
    }

    const parsed = parseScriptureReference(trimmed, books);

    if (parsed === null) {
        activeRequest++;
        crossrefs.loading = false;
        crossrefs.references = [];

        return;
    }

    debounceTimeout = setTimeout(() => {
        void loadCrossReferences(parsed.bookId, parsed.chapter, parsed.verse);
    }, 300);
}

export async function loadCrossReferences(bookId: string, chapter: number, verse: number): Promise<void> {
    const requestId = ++activeRequest;
    crossrefs.loading = true;

    try {
        const url = crossReferencesRoute.url({ query: { book: bookId, chapter, verse } });
        const response = await fetch(url, {
            headers: { Accept: 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
        });

        if (requestId !== activeRequest) {
            return;
        }

        if (! response.ok) {
            crossrefs.references = [];

            return;
        }

        const data = await response.json() as { references: CrossReference[] };

        if (requestId !== activeRequest) {
            return;
        }

        crossrefs.references = data.references;
    } finally {
        if (requestId === activeRequest) {
            crossrefs.loading = false;
        }
    }
}
