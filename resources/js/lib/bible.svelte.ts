import {
    books as booksRoute,
    chapter as chapterRoute,
    translations as translationsRoute,
} from '@/actions/App/Http/Controllers/BibleController';
import type { Book, Chapter, Translation } from '@/lib/types/bible';

export const bible = $state({
    translations: [] as Translation[],
    books: [] as Book[],
    booksTranslationId: null as string | null,
    translationsLoading: false,
    booksLoading: false,
    translationsLoaded: false,
});

const chapterCache = new Map<string, Chapter>();
const chapterInflight = new Map<string, Promise<Chapter>>();

async function jsonFetch<T>(url: string): Promise<T> {
    const response = await fetch(url, {
        headers: {
            Accept: 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
        },
    });

    if (! response.ok) {
        throw new Error(`Request failed: ${response.status}`);
    }

    return response.json() as Promise<T>;
}

export async function loadTranslations(): Promise<void> {
    if (bible.translationsLoaded) {
        return;
    }

    bible.translationsLoading = true;

    try {
        const data = await jsonFetch<{ translations: Translation[] }>(translationsRoute.url());
        bible.translations = data.translations;
        bible.translationsLoaded = true;
    } finally {
        bible.translationsLoading = false;
    }
}

export async function loadBooks(translationId: string): Promise<void> {
    if (bible.booksTranslationId === translationId && bible.books.length > 0) {
        return;
    }

    bible.booksLoading = true;

    try {
        const data = await jsonFetch<{ books: Book[] }>(booksRoute.url(translationId));
        bible.books = data.books;
        bible.booksTranslationId = translationId;
    } finally {
        bible.booksLoading = false;
    }
}

function chapterKey(translationId: string, bookId: string, chapter: number): string {
    return `${translationId}:${bookId}:${chapter}`;
}

export async function fetchChapter(
    translationId: string,
    bookId: string,
    chapterNumber: number,
): Promise<Chapter> {
    const key = chapterKey(translationId, bookId, chapterNumber);
    const cached = chapterCache.get(key);

    if (cached) {
        return cached;
    }

    const inflight = chapterInflight.get(key);

    if (inflight) {
        return inflight;
    }

    const promise = jsonFetch<{ chapter: Chapter }>(
        chapterRoute.url({
            translation: translationId,
            book: bookId,
            chapter: chapterNumber,
        }),
    )
        .then((data) => {
            chapterCache.set(key, data.chapter);
            chapterInflight.delete(key);

            return data.chapter;
        })
        .catch((error) => {
            chapterInflight.delete(key);

            throw error;
        });

    chapterInflight.set(key, promise);

    return promise;
}

export function getAdjacentChapter(
    bookId: string,
    chapter: number,
    direction: 'prev' | 'next',
): { bookId: string; chapter: number } | null {
    const book = bible.books.find((item) => item.id === bookId);

    if (! book) {
        return null;
    }

    if (direction === 'prev') {
        if (chapter > 1) {
            return { bookId, chapter: chapter - 1 };
        }

        const index = bible.books.findIndex((item) => item.id === bookId);

        if (index > 0) {
            const previous = bible.books[index - 1];

            return { bookId: previous.id, chapter: previous.chapters };
        }

        return null;
    }

    if (chapter < book.chapters) {
        return { bookId, chapter: chapter + 1 };
    }

    const index = bible.books.findIndex((item) => item.id === bookId);

    if (index < bible.books.length - 1) {
        const next = bible.books[index + 1];

        return { bookId: next.id, chapter: 1 };
    }

    return null;
}

export function bookAbbrev(bookId: string): string {
    return bible.books.find((item) => item.id === bookId)?.abbrev ?? bookId.toUpperCase();
}
