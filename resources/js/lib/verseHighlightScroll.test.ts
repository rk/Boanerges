import { describe, expect, it, vi } from 'vitest';

import { createVerseHighlightScroller } from '@/lib/verseHighlightScroll';

function scrollRootWithVerse(verse: number): HTMLElement {
    const root = document.createElement('div');
    const verseEl = document.createElement('span');

    verseEl.setAttribute('data-verse', String(verse));
    verseEl.scrollIntoView = vi.fn();
    root.appendChild(verseEl);

    return root;
}

describe('createVerseHighlightScroller', () => {
    it('scrolls every root when all expected columns are ready', async () => {
        const scroller = createVerseHighlightScroller();
        const first = scrollRootWithVerse(5);
        const second = scrollRootWithVerse(5);

        await scroller.scrollTo(
            [first, second],
            { verse: 5, endVerse: null },
            'gen:1',
            2,
        );

        expect(
            first.querySelector('[data-verse="5"]')?.scrollIntoView,
        ).toHaveBeenCalled();
        expect(
            second.querySelector('[data-verse="5"]')?.scrollIntoView,
        ).toHaveBeenCalled();
    });

    it('retries until verse elements exist in every root', async () => {
        const scroller = createVerseHighlightScroller();
        const root = document.createElement('div');

        const pending = scroller.scrollTo(
            root,
            { verse: 3, endVerse: null },
            'gen:1',
            1,
        );

        await new Promise((resolve) => setTimeout(resolve, 80));

        const verseEl = document.createElement('span');
        verseEl.setAttribute('data-verse', '3');
        verseEl.scrollIntoView = vi.fn();
        root.appendChild(verseEl);

        await pending;

        expect(verseEl.scrollIntoView).toHaveBeenCalled();
    });

    it('skips duplicate scroll requests for the same highlight key', async () => {
        const scroller = createVerseHighlightScroller();
        const root = scrollRootWithVerse(3);
        const verseEl = root.querySelector(
            '[data-verse="3"]',
        ) as HTMLElement & {
            scrollIntoView: ReturnType<typeof vi.fn>;
        };

        await scroller.scrollTo(root, { verse: 3, endVerse: null }, 'gen:1', 1);
        await scroller.scrollTo(root, { verse: 3, endVerse: null }, 'gen:1', 1);

        expect(verseEl.scrollIntoView).toHaveBeenCalledTimes(1);
    });
});
