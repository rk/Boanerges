const STORAGE_PREFIX = 'boanerges.scribe';

function storageKey(bookId: string, chapter: number): string {
    return `${STORAGE_PREFIX}.${bookId}.${chapter}`;
}

export const scribe = $state({
    saveStatus: 'idle' as 'idle' | 'saving' | 'saved',
});

let saveTimeout: ReturnType<typeof setTimeout> | null = null;
let hideTimeout: ReturnType<typeof setTimeout> | null = null;

export function loadScribeDraft(bookId: string, chapterNumber: number): Record<number, string> {
    if (typeof window === 'undefined') {
        return {};
    }

    try {
        const stored = localStorage.getItem(storageKey(bookId, chapterNumber));

        return stored ? JSON.parse(stored) : {};
    } catch {
        return {};
    }
}

export function scheduleScribeSave(
    bookId: string,
    chapterNumber: number,
    draft: Record<number, string>,
): void {
    scribe.saveStatus = 'saving';

    if (saveTimeout) {
        clearTimeout(saveTimeout);
    }

    if (hideTimeout) {
        clearTimeout(hideTimeout);
    }

    saveTimeout = setTimeout(() => {
        if (typeof window !== 'undefined') {
            localStorage.setItem(storageKey(bookId, chapterNumber), JSON.stringify(draft));
        }

        scribe.saveStatus = 'saved';

        hideTimeout = setTimeout(() => {
            scribe.saveStatus = 'idle';
        }, 1500);
    }, 500);
}
