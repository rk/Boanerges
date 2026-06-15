import { SvelteMap } from 'svelte/reactivity';
import {
    books as booksRoute,
    catalog as catalogRoute,
    chapter as chapterRoute,
    install as installRoute,
    translations as translationsRoute,
    uninstall as uninstallRoute,
} from '@/actions/App/Http/Controllers/BibleController';
import type { Book, CatalogTranslation, Chapter, Translation } from '@/lib/types/bible';

export const bible = $state({
    translations: [] as Translation[],
    catalog: [] as CatalogTranslation[],
    books: [] as Book[],
    booksTranslationId: null as string | null,
    translationsLoading: false,
    catalogLoading: false,
    booksLoading: false,
    translationsLoaded: false,
    catalogLoaded: false,
    translationManagerOpen: false,
    installingModule: null as string | null,
    uninstallingModule: null as string | null,
    managerError: null as string | null,
});

const chapterCache = new SvelteMap<string, Chapter>();
const chapterInflight = new SvelteMap<string, Promise<Chapter>>();

async function jsonFetch<T>(url: string, init?: RequestInit): Promise<T> {
    const response = await fetch(url, {
        ...init,
        headers: {
            Accept: 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
            'Content-Type': 'application/json',
            ...(init?.headers ?? {}),
        },
    });

    if (! response.ok) {
        const body = (await response.json().catch(() => null)) as { message?: string } | null;

        throw new Error(body?.message ?? `Request failed: ${response.status}`);
    }

    return response.json() as Promise<T>;
}

export function openTranslationManager(): void {
    bible.translationManagerOpen = true;
    bible.managerError = null;
    void loadCatalog();
}

export function closeTranslationManager(): void {
    bible.translationManagerOpen = false;
    bible.managerError = null;
}

export function invalidateTranslations(): void {
    bible.translationsLoaded = false;
    bible.catalogLoaded = false;
    bible.translations = [];
    bible.catalog = [];
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

export async function loadCatalog(force = false): Promise<void> {
    if (bible.catalogLoaded && ! force) {
        return;
    }

    bible.catalogLoading = true;

    try {
        const data = await jsonFetch<{ translations: CatalogTranslation[] }>(catalogRoute.url());
        bible.catalog = data.translations;
        bible.catalogLoaded = true;
    } finally {
        bible.catalogLoading = false;
    }
}

export async function installTranslation(module: string): Promise<void> {
    bible.installingModule = module;
    bible.managerError = null;

    try {
        await jsonFetch(installRoute.url(module), { method: 'POST' });
        invalidateTranslations();
        await Promise.all([loadTranslations(), loadCatalog(true)]);
    } catch (error) {
        bible.managerError = error instanceof Error ? error.message : 'Installation failed.';

        throw error;
    } finally {
        bible.installingModule = null;
    }
}

export async function uninstallTranslation(module: string): Promise<void> {
    bible.uninstallingModule = module;
    bible.managerError = null;

    try {
        await jsonFetch(uninstallRoute.url(module), { method: 'DELETE' });
        invalidateTranslations();
        await Promise.all([loadTranslations(), loadCatalog(true)]);
    } catch (error) {
        bible.managerError = error instanceof Error ? error.message : 'Removal failed.';

        throw error;
    } finally {
        bible.uninstallingModule = null;
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
