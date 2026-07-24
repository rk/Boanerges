import type { SearchResult } from '@/lib/types/bible';
import { search as searchRoute } from '@/actions/App/Http/Controllers/BibleController';

export const search = $state({
    loading: false,
    results: [] as SearchResult[],
    hasMore: false,
    offset: 0,
    query: '',
});

let debounceTimeout: ReturnType<typeof setTimeout> | null = null;
let activeRequest = 0;

const PAGE_SIZE = 25;

export function resetSearchState(): void {
    search.results = [];
    search.hasMore = false;
    search.offset = 0;
}

export function scheduleSearch(query: string, translationId: string): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
    }

    const trimmed = query.trim();
    search.query = trimmed;

    if (trimmed.length < 2) {
        activeRequest++;
        search.loading = false;
        resetSearchState();

        return;
    }

    debounceTimeout = setTimeout(() => {
        resetSearchState();
        void executeSearch(trimmed, translationId, 0, false);
    }, 300);
}

export function cancelSearch(): void {
    if (debounceTimeout) {
        clearTimeout(debounceTimeout);
        debounceTimeout = null;
    }

    activeRequest++;
    search.loading = false;
    search.query = '';
    resetSearchState();
}

export function loadMoreSearchResults(translationId: string): void {
    if (!search.hasMore || search.loading || search.query.length < 2) {
        return;
    }

    void executeSearch(search.query, translationId, search.offset, true);
}

async function executeSearch(
    query: string,
    translationId: string,
    offset: number,
    append: boolean,
): Promise<void> {
    const requestId = ++activeRequest;
    search.loading = true;

    try {
        const url = searchRoute.url({
            query: {
                q: query,
                translation: translationId,
                limit: PAGE_SIZE,
                offset,
            },
        });
        const response = await fetch(url, {
            headers: {
                Accept: 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
            },
        });

        if (requestId !== activeRequest) {
            return;
        }

        if (!response.ok) {
            if (!append) {
                resetSearchState();
            }

            return;
        }

        const data = (await response.json()) as {
            results: SearchResult[];
            hasMore: boolean;
        };

        if (requestId !== activeRequest) {
            return;
        }

        search.results = append
            ? [...search.results, ...data.results]
            : data.results;
        search.hasMore = data.hasMore;
        search.offset = offset + data.results.length;
    } finally {
        if (requestId === activeRequest) {
            search.loading = false;
        }
    }
}
