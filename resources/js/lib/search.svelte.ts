import { search as searchRoute } from '@/actions/App/Http/Controllers/BibleController';
import type { SearchResult } from '@/lib/types/bible';

export const search = $state({
    open: false,
    loading: false,
    results: [] as SearchResult[],
});

let debounceTimeout: ReturnType<typeof setTimeout> | null = null;
let activeRequest = 0;

export function scheduleSearch(query: string, translationId: string): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
    }

    const trimmed = query.trim();

    if (trimmed.length < 2) {
        activeRequest++;
        search.loading = false;
        search.results = [];

        return;
    }

    debounceTimeout = setTimeout(() => {
        void executeSearch(trimmed, translationId);
    }, 300);
}

export function cancelSearch(): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
        debounceTimeout = null;
    }

    activeRequest++;
    search.loading = false;
    search.results = [];
}

async function executeSearch(query: string, translationId: string): Promise<void> {
    const requestId = ++activeRequest;
    search.loading = true;

    try {
        const url = searchRoute.url({ query: { q: query, translation: translationId } });
        const response = await fetch(url, {
            headers: { Accept: 'application/json', 'X-Requested-With': 'XMLHttpRequest' },
        });

        if (requestId !== activeRequest) {
            return;
        }

        if (! response.ok) {
            search.results = [];

            return;
        }

        const data = await response.json() as { results: SearchResult[] };

        if (requestId !== activeRequest) {
            return;
        }

        search.results = data.results;
    } finally {
        if (requestId === activeRequest) {
            search.loading = false;
        }
    }
}

export function openSearch(): void {
    search.open = true;
}
