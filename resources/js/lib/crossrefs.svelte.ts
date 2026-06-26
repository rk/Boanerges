import { crossReferences as crossReferencesRoute } from '@/actions/App/Http/Controllers/BibleController';
import { parseScriptureReference } from '@/lib/scriptureReference';
import type { Book, CrossReference } from '@/lib/types/bible';

export const crossrefs = $state({
    loading: false,
    references: [] as CrossReference[],
    activeReference: null as string | null,
});

let debounceTimeout: ReturnType<typeof setTimeout> | null = null;
let activeRequest = 0;

export function setCrossReferenceInput(reference: string): void {
    crossrefs.activeReference = reference;
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
