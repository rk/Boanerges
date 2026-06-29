import { show, update } from '@/actions/App/Http/Controllers/NotesController';

export const notes = $state({
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

    if (!response.ok) {
        const body = (await response.json().catch(() => null)) as {
            message?: string;
        } | null;

        throw new Error(body?.message ?? `Request failed: ${response.status}`);
    }

    return response.json() as Promise<T>;
}

export async function fetchNotesDraft(
    bookId: string,
    chapter: number,
): Promise<string> {
    const data = await jsonFetch<{ content: string }>(
        show.url({ book: bookId, chapter }),
    );

    return data.content;
}

export function scheduleNotesSave(
    bookId: string,
    chapter: number,
    content: string,
): void {
    if (saveTimeout) {
        clearTimeout(saveTimeout);
    }

    notes.saveStatus = 'idle';

    saveTimeout = setTimeout(() => {
        void persistNotesDraft(bookId, chapter, content);
    }, 500);
}

async function persistNotesDraft(
    bookId: string,
    chapter: number,
    content: string,
): Promise<void> {
    const generation = ++saveGeneration;
    notes.saveStatus = 'saving';

    try {
        await jsonFetch<{ content: string }>(
            update.url({ book: bookId, chapter }),
            {
                method: 'PUT',
                body: JSON.stringify({ content }),
            },
        );

        if (generation !== saveGeneration) {
            return;
        }

        notes.saveStatus = 'saved';

        if (hideTimeout) {
            clearTimeout(hideTimeout);
        }

        hideTimeout = setTimeout(() => {
            if (generation === saveGeneration) {
                notes.saveStatus = 'idle';
            }
        }, 1500);
    } catch {
        if (generation === saveGeneration) {
            notes.saveStatus = 'idle';
        }
    }
}
