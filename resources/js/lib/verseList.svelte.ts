import { SvelteMap } from 'svelte/reactivity';
import { bible, fetchChapter } from '@/lib/bible.svelte.ts';
import { formatScriptureReference } from '@/lib/scriptureReference';
import type { ScriptureReference } from '@/lib/scriptureReference';
import type { Book } from '@/lib/types/bible';
import { verseListNeedsTextLoad as needsTextLoad } from '@/lib/verseListTextLoad';
import {
    destroy as destroyVerseListRoute,
    index as verseListsRoute,
    show as showVerseListRoute,
    store as storeVerseListRoute,
    update as updateVerseListRoute,
} from '@/actions/App/Http/Controllers/VerseListController';

export type VerseListEntry = ScriptureReference & {
    label: string;
    text?: string;
};

export type VerseListSummary = {
    id: string;
    title: string;
    updatedAt: string;
};

export const verseList = $state({
    title: 'Untitled list',
    savedListId: null as string | null,
    showContent: true,
    loading: false,
    saving: false,
    entries: [] as VerseListEntry[],
    savedLists: [] as VerseListSummary[],
});

let textRequest = 0;
let loadedTextTranslationId: string | null = null;

export function verseListNeedsTextLoad(
    entries: ReadonlyArray<{ text?: string }>,
    showContent: boolean,
    translationId: string,
    loading: boolean,
): boolean {
    return needsTextLoad(
        entries,
        showContent,
        translationId,
        loading,
        loadedTextTranslationId,
    );
}

export function hydrateVerseListSettings(
    showContent: boolean,
    activeId: string | null,
): void {
    verseList.showContent = showContent;

    if (activeId) {
        void loadSavedVerseList(activeId);
    }
}

export function setVerseListShowContent(enabled: boolean): void {
    verseList.showContent = enabled;

    if (!enabled) {
        loadedTextTranslationId = null;
        verseList.entries = verseList.entries.map((entry) => ({
            ...entry,
            text: undefined,
        }));
    }
}

export function addVerseToList(ref: ScriptureReference, books: Book[]): void {
    const exists = verseList.entries.some(
        (entry) =>
            entry.bookId === ref.bookId &&
            entry.chapter === ref.chapter &&
            entry.verse === ref.verse,
    );

    if (exists) {
        return;
    }

    verseList.entries = [
        ...verseList.entries,
        {
            ...ref,
            label: formatScriptureReference(
                ref.bookId,
                ref.chapter,
                ref.verse,
                books,
            ),
        },
    ];
}

export function removeVerseFromList(index: number): void {
    verseList.entries = verseList.entries.filter((_, i) => i !== index);
}

export async function loadSavedLists(): Promise<void> {
    const response = await fetch(verseListsRoute.url(), {
        headers: {
            Accept: 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
        },
    });

    if (!response.ok) {
        return;
    }

    const data = (await response.json()) as { lists: VerseListSummary[] };
    verseList.savedLists = data.lists;
}

export async function loadSavedVerseList(id: string): Promise<void> {
    verseList.loading = true;

    try {
        const response = await fetch(showVerseListRoute.url(id), {
            headers: {
                Accept: 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
            },
        });

        if (!response.ok) {
            return;
        }

        const data = (await response.json()) as {
            list: {
                id: string;
                title: string;
                entries: ScriptureReference[];
            };
        };

        verseList.savedListId = data.list.id;
        verseList.title = data.list.title;
        verseList.entries = data.list.entries.map((entry) => ({
            ...entry,
            label: formatScriptureReference(
                entry.bookId,
                entry.chapter,
                entry.verse,
                bible.books,
            ),
        }));
    } finally {
        verseList.loading = false;
    }
}

export async function saveVerseList(): Promise<boolean> {
    verseList.saving = true;

    try {
        const payload = {
            title: verseList.title,
            entries: verseList.entries.map(({ bookId, chapter, verse }) => ({
                bookId,
                chapter,
                verse,
            })),
        };

        const url = verseList.savedListId
            ? updateVerseListRoute.url(verseList.savedListId)
            : storeVerseListRoute.url();

        const response = await fetch(url, {
            method: 'PUT',
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
                'X-Requested-With': 'XMLHttpRequest',
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            return false;
        }

        const data = (await response.json()) as {
            list: { id: string; title: string };
        };

        verseList.savedListId = data.list.id;
        verseList.title = data.list.title;
        await loadSavedLists();

        return true;
    } finally {
        verseList.saving = false;
    }
}

export async function deleteSavedVerseList(id: string): Promise<void> {
    await fetch(destroyVerseListRoute.url(id), {
        method: 'DELETE',
        headers: {
            Accept: 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
        },
    });

    if (verseList.savedListId === id) {
        verseList.savedListId = null;
    }

    await loadSavedLists();
}

export async function loadVerseTexts(translationId: string): Promise<void> {
    if (
        !verseListNeedsTextLoad(
            verseList.entries,
            verseList.showContent,
            translationId,
            verseList.loading,
        )
    ) {
        return;
    }

    const requestId = ++textRequest;

    verseList.loading = true;

    try {
        const chapterCache = new SvelteMap<
            string,
            Awaited<ReturnType<typeof fetchChapter>>
        >();

        const loaded = await Promise.all(
            verseList.entries.map(async (entry) => {
                const key = `${entry.bookId}:${entry.chapter}`;

                if (!chapterCache.has(key)) {
                    try {
                        chapterCache.set(
                            key,
                            await fetchChapter(
                                translationId,
                                entry.bookId,
                                entry.chapter,
                            ),
                        );
                    } catch {
                        chapterCache.set(key, { verses: [] } as never);
                    }
                }

                const chapter = chapterCache.get(key);
                const verse = chapter?.verses.find(
                    (row) => row.number === entry.verse,
                );

                return {
                    bookId: entry.bookId,
                    chapter: entry.chapter,
                    verse: entry.verse,
                    text: verse?.text ?? '—',
                };
            }),
        );

        if (requestId !== textRequest) {
            return;
        }

        loadedTextTranslationId = translationId;

        for (const row of loaded) {
            const entry = verseList.entries.find(
                (candidate) =>
                    candidate.bookId === row.bookId &&
                    candidate.chapter === row.chapter &&
                    candidate.verse === row.verse,
            );

            if (entry) {
                entry.text = row.text;
            }
        }
    } finally {
        if (requestId === textRequest) {
            verseList.loading = false;
        }
    }
}

export function verseListPrintPayload(): {
    title: string;
    entries: ScriptureReference[];
} {
    return {
        title: verseList.title,
        entries: verseList.entries.map(({ bookId, chapter, verse }) => ({
            bookId,
            chapter,
            verse,
        })),
    };
}
