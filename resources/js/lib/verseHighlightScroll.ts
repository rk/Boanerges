import { tick } from 'svelte';

import type { VerseHighlight } from '@/lib/verseHighlight';

const RETRY_MS = 50;
const TIMEOUT_MS = 3000;

function sleep(ms: number): Promise<void> {
    return new Promise((resolve) => setTimeout(resolve, ms));
}

function verseSelector(verse: number): string {
    return `[data-verse="${verse}"]`;
}

export function createVerseHighlightScroller(): {
    reset: () => void;
    scrollTo: (
        scrollRoot: HTMLElement | null | undefined | Array<HTMLElement | null | undefined>,
        highlight: VerseHighlight,
        locationKey: string,
        expectedRootCount?: number,
        isCancelled?: () => boolean,
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
            expectedRootCount = 1,
            isCancelled?: () => boolean,
        ): Promise<void> {
            const key = `${locationKey}:${highlight.verse}:${highlight.endVerse ?? ''}`;

            if (key === lastKey) {
                return;
            }

            const deadline = Date.now() + TIMEOUT_MS;
            const behavior = expectedRootCount > 1 ? 'instant' : 'smooth';

            while (Date.now() < deadline) {
                if (isCancelled?.()) {
                    return;
                }

                await tick();
                await new Promise((resolve) => requestAnimationFrame(resolve));

                const roots = (Array.isArray(scrollRoot) ? scrollRoot : [scrollRoot])
                    .filter((root): root is HTMLElement => root instanceof HTMLElement);

                if (roots.length < expectedRootCount) {
                    await sleep(RETRY_MS);

                    continue;
                }

                const verseElements = roots.map((root) =>
                    root.querySelector(verseSelector(highlight.verse)),
                );

                if (verseElements.some((element) => element === null)) {
                    await sleep(RETRY_MS);

                    continue;
                }

                for (const element of verseElements) {
                    element?.scrollIntoView({ behavior, block: 'center' });
                }

                lastKey = key;

                return;
            }
        },
    };
}

export function highlightedVersesFromHighlight(
    highlight: VerseHighlight | null,
    verseNumbersInRange: (verse: number, endVerse: number | null) => Set<number>,
): Set<number> {
    return highlight ? verseNumbersInRange(highlight.verse, highlight.endVerse) : new Set<number>();
}
