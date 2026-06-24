import { tick } from 'svelte';

import type { VerseHighlight } from '@/lib/verseHighlight';

export function createVerseHighlightScroller(): {
    reset: () => void;
    scrollTo: (
        scrollRoot: HTMLElement | null | undefined | Array<HTMLElement | null | undefined>,
        highlight: VerseHighlight,
        locationKey: string,
    ) => Promise<void>;
} {
    let lastKey: string | null = null;

    return {
        reset(): void {
            lastKey = null;
        },
        async scrollTo(
            scrollRoot: HTMLElement | null | undefined | Array<HTMLElement | null | undefined>,
            highlight: VerseHighlight,
            locationKey: string,
        ): Promise<void> {
            const key = `${locationKey}:${highlight.verse}:${highlight.endVerse ?? ''}`;

            if (key === lastKey) {
                return;
            }

            await tick();

            const roots = Array.isArray(scrollRoot) ? scrollRoot : [scrollRoot];
            const behavior = roots.length > 1 ? 'instant' : 'smooth';

            for (const root of roots) {
                root?.querySelector(`[data-verse="${highlight.verse}"]`)
                    ?.scrollIntoView({ behavior, block: 'center' });
            }

            lastKey = key;
        },
    };
}

export function highlightedVersesFromHighlight(
    highlight: VerseHighlight | null,
    verseNumbersInRange: (verse: number, endVerse: number | null) => Set<number>,
): Set<number> {
    return highlight ? verseNumbersInRange(highlight.verse, highlight.endVerse) : new Set<number>();
}
