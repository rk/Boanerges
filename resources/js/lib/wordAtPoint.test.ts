import { afterEach, describe, expect, it, vi } from 'vitest';

import {
    extractWordFromTextOffset,
    getLastWordLookupMethod,
    wordAtPoint,
} from '@/lib/wordAtPoint';

describe('extractWordFromTextOffset', () => {
    it('extracts the word surrounding an offset', () => {
        expect(extractWordFromTextOffset('the grace of God', 8)).toBe('grace');
    });

    it('returns empty string when offset is outside a word', () => {
        expect(extractWordFromTextOffset('the   grace', 4)).toBe('');
    });
});

describe('wordAtPoint', () => {
    const originalCaretPositionFromPoint =
        document.caretPositionFromPoint?.bind(document);
    const originalCaretRangeFromPoint = (
        document as Document & {
            caretRangeFromPoint?: (x: number, y: number) => Range | null;
        }
    ).caretRangeFromPoint?.bind(document);

    afterEach(() => {
        document.body.innerHTML = '';
        vi.restoreAllMocks();

        if (originalCaretPositionFromPoint) {
            document.caretPositionFromPoint = originalCaretPositionFromPoint;
        }

        if (originalCaretRangeFromPoint) {
            (
                document as Document & {
                    caretRangeFromPoint?: (
                        x: number,
                        y: number,
                    ) => Range | null;
                }
            ).caretRangeFromPoint = originalCaretRangeFromPoint;
        }
    });

    it('prefers caretRangeFromPoint when available', () => {
        document.body.innerHTML =
            '<p data-verse="1"><span data-dict-word="grace">grace</span></p>';
        const textNode = document.querySelector('[data-dict-word]')!
            .firstChild as Text;
        const target = document.querySelector('[data-dict-word]')!;
        const range = document.createRange();
        range.setStart(textNode, 3);
        range.collapse(true);

        (
            document as Document & {
                caretRangeFromPoint: (x: number, y: number) => Range | null;
            }
        ).caretRangeFromPoint = () => range;

        const event = new MouseEvent('contextmenu', {
            bubbles: true,
            clientX: 10,
            clientY: 10,
        });
        Object.defineProperty(event, 'target', { value: target });

        const debugSpy = vi
            .spyOn(console, 'debug')
            .mockImplementation(() => {});

        expect(wordAtPoint(event)).toBe('grace');
        expect(getLastWordLookupMethod()).toBe('caretRangeFromPoint');
        expect(debugSpy).toHaveBeenCalledWith('[wordAtPoint]', {
            method: 'caretRangeFromPoint',
            word: 'grace',
        });
    });

    it('falls back to data-dict-word when caret APIs miss', () => {
        document.body.innerHTML =
            '<p data-verse="1"><span data-dict-word="grace">grace</span></p>';

        (
            document as Document & {
                caretRangeFromPoint: () => Range | null;
            }
        ).caretRangeFromPoint = () => null;
        document.caretPositionFromPoint = () => null;

        const target = document.querySelector('[data-dict-word]')!;
        const event = new MouseEvent('contextmenu', {
            bubbles: true,
            clientX: 10,
            clientY: 10,
        });
        Object.defineProperty(event, 'target', { value: target });
        const debugSpy = vi
            .spyOn(console, 'debug')
            .mockImplementation(() => {});

        expect(wordAtPoint(event)).toBe('grace');
        expect(getLastWordLookupMethod()).toBe('data-dict-word');
        expect(debugSpy).toHaveBeenCalledWith('[wordAtPoint]', {
            method: 'data-dict-word',
            word: 'grace',
        });
    });

    it('uses caretPositionFromPoint when caretRangeFromPoint misses', () => {
        document.body.innerHTML =
            '<p data-verse="1"><span>the grace of God</span></p>';
        const textNode = document.querySelector('span')!.firstChild as Text;

        (
            document as Document & {
                caretRangeFromPoint: () => Range | null;
            }
        ).caretRangeFromPoint = () => null;
        document.caretPositionFromPoint = () => ({
            offsetNode: textNode,
            offset: 8,
            getClientRect: () => new DOMRect(),
        });

        const target = document.querySelector('span')!;
        const event = new MouseEvent('contextmenu', {
            bubbles: true,
            clientX: 10,
            clientY: 10,
        });
        Object.defineProperty(event, 'target', { value: target });
        const debugSpy = vi
            .spyOn(console, 'debug')
            .mockImplementation(() => {});

        expect(wordAtPoint(event)).toBe('grace');
        expect(getLastWordLookupMethod()).toBe('caretPositionFromPoint');
        expect(debugSpy).toHaveBeenCalledWith('[wordAtPoint]', {
            method: 'caretPositionFromPoint',
            word: 'grace',
        });
    });
});
