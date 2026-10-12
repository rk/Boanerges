import { describe, expect, it } from 'vitest';

import {
    applyParagraphStart,
    focusAdjacentVerse,
    hydrateDocument,
    hydrateVerseSpan,
    normalizePasteText,
    paragraphStartAttribute,
    readVerseSpan,
    serializeFromDocument,
    shouldHydrateVerseSpan,
} from '@/lib/scribeDocument';

describe('paragraphStartAttribute', () => {
    it('maps booleans to data attribute values', () => {
        expect(paragraphStartAttribute(true)).toBe('true');
        expect(paragraphStartAttribute(false)).toBe('false');
    });
});

describe('normalizePasteText', () => {
    it('normalizes Windows and classic Mac line endings', () => {
        expect(normalizePasteText('a\r\nb\rc')).toBe('a\nb\nc');
    });
});

describe('verse span hydration', () => {
    it('hydrates text when the span is not focused', () => {
        const span = document.createElement('span');

        hydrateVerseSpan(span, 'In the beginning');

        expect(readVerseSpan(span)).toBe('In the beginning');
    });

    it('skips hydration while focused unless forced', () => {
        const span = document.createElement('span');
        span.textContent = 'live draft';
        document.body.appendChild(span);
        span.focus();

        expect(shouldHydrateVerseSpan(span, 'server draft')).toBe(false);

        hydrateVerseSpan(span, 'server draft');

        expect(readVerseSpan(span)).toBe('live draft');

        hydrateVerseSpan(span, 'server draft', { force: true });

        expect(readVerseSpan(span)).toBe('server draft');
    });

    it('applies paragraph start metadata on wrappers', () => {
        const span = document.createElement('span');

        applyParagraphStart(span, true);

        expect(span.dataset.paragraphStart).toBe('true');
    });
});

describe('document serialization', () => {
    it('round-trips verse text and preserves intra-verse blank lines', () => {
        const root = document.createElement('div');
        root.innerHTML = `
            <span class="scribe-verse-wrap" data-paragraph-start="true">
                <span class="scribe-verse" data-verse="1">First line\n\nSecond line</span>
            </span>
            <span class="scribe-verse-wrap" data-paragraph-start="false">
                <span class="scribe-verse" data-verse="2"></span>
            </span>
            <span class="scribe-verse-wrap" data-paragraph-start="true">
                <span class="scribe-verse" data-verse="3">Third verse</span>
            </span>
        `;

        expect(serializeFromDocument(root)).toEqual({
            1: { text: 'First line\n\nSecond line' },
            2: { text: '' },
            3: { text: 'Third verse' },
        });
    });

    it('hydrates all verse spans from entries', () => {
        const root = document.createElement('div');
        root.innerHTML = `
            <span class="scribe-verse" data-verse="1"></span>
            <span class="scribe-verse" data-verse="2"></span>
        `;

        hydrateDocument(
            root,
            {
                1: { text: 'Alpha' },
                2: { text: 'Beta' },
            },
            {
                1: true,
                2: false,
            },
            { force: true },
        );

        const spans = [...root.querySelectorAll<HTMLElement>('.scribe-verse')];

        expect(readVerseSpan(spans[0])).toBe('Alpha');
        expect(readVerseSpan(spans[1])).toBe('Beta');
        expect(spans[0].dataset.paragraphStart).toBe('true');
        expect(spans[1].dataset.paragraphStart).toBe('false');
    });
});

describe('focusAdjacentVerse', () => {
    it('moves focus between verse spans with Tab semantics', () => {
        const root = document.createElement('div');
        root.innerHTML = `
            <span class="scribe-verse" data-verse="1" tabindex="0"></span>
            <span class="scribe-verse" data-verse="2" tabindex="0"></span>
        `;

        document.body.appendChild(root);

        const spans = [...root.querySelectorAll<HTMLElement>('.scribe-verse')];
        spans[0].focus();

        expect(focusAdjacentVerse(root, 1, 1)).toBe(true);
        expect(document.activeElement).toBe(spans[1]);
    });
});
