import { show, update } from '@/actions/App/Http/Controllers/ScribeController';
import type { ScribeVerse } from '@/lib/types/bible';

const LEGACY_STORAGE_PREFIX = 'boanerges.scribe';

export const scribe = $state({
    saveStatus: 'idle' as 'idle' | 'saving' | 'saved',
});

let saveTimeout: ReturnType<typeof setTimeout> | null = null;
let hideTimeout: ReturnType<typeof setTimeout> | null = null;
let saveGeneration = 0;

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

function legacyStorageKey(bookId: string, chapter: number): string {
    return `${LEGACY_STORAGE_PREFIX}.${bookId}.${chapter}`;
}

function loadLegacyDraft(bookId: string, chapter: number): ScribeVerse[] {
    if (typeof window === 'undefined') {
        return [];
    }

    try {
        const stored = localStorage.getItem(legacyStorageKey(bookId, chapter));

        if (! stored) {
            return [];
        }

        const parsed = JSON.parse(stored) as Record<string, string>;

        return Object.entries(parsed).map(([verse, text]) => ({
            verse: Number(verse),
            text,
        }));
    } catch {
        return [];
    }
}

function clearLegacyDraft(bookId: string, chapter: number): void {
    if (typeof window === 'undefined') {
        return;
    }

    localStorage.removeItem(legacyStorageKey(bookId, chapter));
}

export async function fetchScribeDraft(bookId: string, chapter: number): Promise<ScribeVerse[]> {
    const data = await jsonFetch<{ verses: ScribeVerse[] }>(show.url({ book: bookId, chapter }));

    if (data.verses.length > 0) {
        return data.verses;
    }

    const legacy = loadLegacyDraft(bookId, chapter);

    if (legacy.length > 0) {
        await jsonFetch<{ verses: ScribeVerse[] }>(update.url({ book: bookId, chapter }), {
            method: 'PUT',
            body: JSON.stringify({ verses: legacy }),
        });
        clearLegacyDraft(bookId, chapter);

        return legacy;
    }

    return [];
}

export function scheduleScribeSave(bookId: string, chapter: number, verses: ScribeVerse[]): void {
    scribe.saveStatus = 'saving';

    if (saveTimeout) {
        clearTimeout(saveTimeout);
    }

    if (hideTimeout) {
        clearTimeout(hideTimeout);
    }

    const generation = ++saveGeneration;

    saveTimeout = setTimeout(() => {
        void jsonFetch<{ verses: ScribeVerse[] }>(update.url({ book: bookId, chapter }), {
            method: 'PUT',
            body: JSON.stringify({ verses }),
        })
            .then(() => {
                if (generation !== saveGeneration) {
                    return;
                }

                scribe.saveStatus = 'saved';

                hideTimeout = setTimeout(() => {
                    if (generation === saveGeneration) {
                        scribe.saveStatus = 'idle';
                    }
                }, 1500);
            })
            .catch(() => {
                if (generation === saveGeneration) {
                    scribe.saveStatus = 'idle';
                }
            });
    }, 500);
}

export function effectiveParagraphStart(
    verseNumber: number,
    sourceParagraphStart?: boolean,
    override?: boolean,
): boolean {
    if (override !== undefined) {
        return override;
    }

    if (sourceParagraphStart) {
        return true;
    }

    return verseNumber === 1;
}

export type ScribeDraftEntry = {
    text: string;
    paragraphStartOverride?: boolean;
};

export function entriesFromScribeVerses(verses: ScribeVerse[]): Record<number, ScribeDraftEntry> {
    const entries: Record<number, ScribeDraftEntry> = {};

    for (const verse of verses) {
        entries[verse.verse] = {
            text: verse.text,
            ...(verse.paragraphStart !== undefined
                ? { paragraphStartOverride: verse.paragraphStart }
                : {}),
        };
    }

    return entries;
}

export function serializeScribeDraft(
    verseNumbers: number[],
    entries: Record<number, ScribeDraftEntry>,
): ScribeVerse[] {
    const verses: ScribeVerse[] = [];

    for (const verseNumber of verseNumbers) {
        const entry = entries[verseNumber];
        const text = entry?.text ?? '';

        if (text.trim() === '' && entry?.paragraphStartOverride === undefined) {
            continue;
        }

        const serialized: ScribeVerse = {
            verse: verseNumber,
            text,
        };

        if (entry?.paragraphStartOverride !== undefined) {
            serialized.paragraphStart = entry.paragraphStartOverride;
        }

        verses.push(serialized);
    }

    return verses;
}
