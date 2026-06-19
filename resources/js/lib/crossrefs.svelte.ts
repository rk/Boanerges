import { crossReferences as crossReferencesRoute } from '@/actions/App/Http/Controllers/BibleController';
import type { CrossReference } from '@/lib/types/bible';

export const crossrefs = $state({
    open: false,
    loading: false,
    references: [] as CrossReference[],
});

export async function loadCrossReferences(bookId: string, chapter: number, verse: number): Promise<void> {
    crossrefs.loading = true;

    try {
        const url = crossReferencesRoute.url({ query: { book: bookId, chapter, verse } });
        const response = await fetch(url, {
            headers: { Accept: 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
        });

        if (! response.ok) {
            crossrefs.references = [];

            return;
        }

        const data = await response.json() as { references: CrossReference[] };
        crossrefs.references = data.references;
    } finally {
        crossrefs.loading = false;
    }
}

export function openCrossReferences(verse: number): void {
    crossrefs.open = true;
}
