import type { ScribeDraftEntry } from '@/lib/scribe.svelte.ts';

export const SCRIBE_VERSE_SELECTOR = '[data-verse]';

export function paragraphStartAttribute(starts: boolean): 'true' | 'false' {
    return starts ? 'true' : 'false';
}

export function normalizePasteText(text: string): string {
    return text.replace(/\r\n/g, '\n').replace(/\r/g, '\n');
}

export function readVerseSpan(element: HTMLElement): string {
    return element.textContent ?? '';
}

export function shouldHydrateVerseSpan(
    element: HTMLElement,
    nextText: string,
    force = false,
): boolean {
    if (force) {
        return true;
    }

    if (document.activeElement === element) {
        return false;
    }

    return readVerseSpan(element) !== nextText;
}

export function hydrateVerseSpan(
    element: HTMLElement,
    text: string,
    options: { force?: boolean } = {},
): void {
    if (!shouldHydrateVerseSpan(element, text, options.force)) {
        return;
    }

    element.textContent = text;
}

export function applyParagraphStart(
    element: HTMLElement,
    starts: boolean,
): void {
    element.dataset.paragraphStart = paragraphStartAttribute(starts);
}

export function verseNumberFromElement(element: HTMLElement): number | null {
    const raw = element.dataset.verse;

    if (raw === undefined) {
        return null;
    }

    const parsed = Number.parseInt(raw, 10);

    return Number.isNaN(parsed) ? null : parsed;
}

export function serializeFromDocument(
    root: HTMLElement,
): Record<number, ScribeDraftEntry> {
    const entries: Record<number, ScribeDraftEntry> = {};

    for (const element of root.querySelectorAll<HTMLElement>(
        SCRIBE_VERSE_SELECTOR,
    )) {
        const verseNumber = verseNumberFromElement(element);

        if (verseNumber === null) {
            continue;
        }

        entries[verseNumber] = {
            text: readVerseSpan(element),
        };
    }

    return entries;
}

export function hydrateDocument(
    root: HTMLElement,
    entries: Record<number, ScribeDraftEntry>,
    paragraphStarts: Record<number, boolean>,
    options: { force?: boolean } = {},
): void {
    for (const element of root.querySelectorAll<HTMLElement>(
        SCRIBE_VERSE_SELECTOR,
    )) {
        const verseNumber = verseNumberFromElement(element);

        if (verseNumber === null) {
            continue;
        }

        hydrateVerseSpan(element, entries[verseNumber]?.text ?? '', options);
        applyParagraphStart(element, paragraphStarts[verseNumber] ?? false);
    }
}

export function focusAdjacentVerse(
    root: HTMLElement,
    currentVerse: number,
    direction: 1 | -1,
): boolean {
    const spans = [
        ...root.querySelectorAll<HTMLElement>(SCRIBE_VERSE_SELECTOR),
    ];
    const index = spans.findIndex(
        (element) => verseNumberFromElement(element) === currentVerse,
    );

    if (index === -1) {
        return false;
    }

    const next = spans[index + direction];

    if (!next) {
        return false;
    }

    next.focus();

    return true;
}
