import { normalizeDictionaryWord } from '@/lib/normalizeDictionaryWord';
import {
    show,
    suggest,
} from '@/actions/App/Http/Controllers/DictionaryController';

export type DictionaryVariant = {
    partOfSpeech: string | null;
    partOfSpeechExpanded: string;
    definitions: string[];
};

export type DictionaryResult = {
    found: boolean;
    word: string;
    variants?: DictionaryVariant[];
    message?: string;
    status?: 'importing';
};

export const dictionary = $state({
    loading: false,
    importing: false,
    result: null as DictionaryResult | null,
    suggestions: [] as string[],
    activeWord: null as string | null,
});

let lookupDebounce: ReturnType<typeof setTimeout> | null = null;
let suggestDebounce: ReturnType<typeof setTimeout> | null = null;
let activeLookupRequest = 0;
let activeSuggestRequest = 0;

export function setDictionaryWord(word: string): void {
    dictionary.activeWord = word;
}

export function scheduleDictionaryLookup(input: string): void {
    if (lookupDebounce) {
        clearTimeout(lookupDebounce);
    }

    const word = normalizeDictionaryWord(input);

    if (word === '') {
        activeLookupRequest++;
        dictionary.loading = false;
        dictionary.importing = false;
        dictionary.result = null;

        return;
    }

    lookupDebounce = setTimeout(() => {
        void loadDictionaryEntry(word);
    }, 300);
}

export function scheduleDictionarySuggest(input: string): void {
    if (suggestDebounce) {
        clearTimeout(suggestDebounce);
    }

    const word = normalizeDictionaryWord(input);

    if (word.length < 2) {
        activeSuggestRequest++;
        dictionary.suggestions = [];

        return;
    }

    suggestDebounce = setTimeout(() => {
        void loadDictionarySuggestions(word);
    }, 200);
}

export async function loadDictionaryEntry(word: string): Promise<void> {
    const requestId = ++activeLookupRequest;
    dictionary.loading = true;
    dictionary.importing = false;

    try {
        const response = await fetch(show.url(word), {
            headers: {
                Accept: 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
            },
        });

        if (requestId !== activeLookupRequest) {
            return;
        }

        if (response.status === 503) {
            dictionary.importing = true;
            dictionary.result = null;

            return;
        }

        if (!response.ok) {
            dictionary.result = null;

            return;
        }

        dictionary.result = (await response.json()) as DictionaryResult;
    } finally {
        if (requestId === activeLookupRequest) {
            dictionary.loading = false;
        }
    }
}

export async function loadDictionarySuggestions(prefix: string): Promise<void> {
    const requestId = ++activeSuggestRequest;

    try {
        const url = suggest.url({
            query: { q: prefix, limit: '10' },
        });
        const response = await fetch(url, {
            headers: {
                Accept: 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
            },
        });

        if (requestId !== activeSuggestRequest) {
            return;
        }

        if (response.status === 503 || !response.ok) {
            dictionary.suggestions = [];

            return;
        }

        const data = (await response.json()) as { suggestions: string[] };
        dictionary.suggestions = data.suggestions;
    } catch {
        if (requestId === activeSuggestRequest) {
            dictionary.suggestions = [];
        }
    }
}
