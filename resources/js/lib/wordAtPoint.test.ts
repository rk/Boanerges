import { afterEach, describe, expect, it } from 'vitest';

import { extractWordFromTextOffset, wordAtPoint } from '@/lib/wordAtPoint';

describe('extractWordFromTextOffset', () => {
    it('extracts the word surrounding an offset', () => {
        expect(extractWordFromTextOffset('the grace of God', 8)).toBe('grace');
    });

    it('returns empty string when offset is outside a word', () => {
        expect(extractWordFromTextOffset('the   grace', 4)).toBe('');
    });
});

describe('wordAtPoint', () => {
    afterEach(() => {
        document.body.innerHTML = '';
        delete (
            document as Document & {
                caretRangeFromPoint?: (x: number, y: number) => Range | null;
                caretPositionFromPoint?: (
                    x: number,
                    y: number,
                ) => { offsetNode: Node; offset: number } | null;
            }
        ).caretRangeFromPoint;
        delete (
            document as Document & {
                caretPositionFromPoint?: (
                    x: number,
                    y: number,
                ) => { offsetNode: Node; offset: number } | null;
            }
        ).caretPositionFromPoint;
    });

    it('reads words from data-dict-word spans', () => {
        document.body.innerHTML =
            '<p data-verse="1"><span data-dict-word="grace">grace</span></p>';

        const target = document.querySelector('[data-dict-word]')!;
        const event = new MouseEvent('contextmenu', {
            bubbles: true,
            clientX: 10,
            clientY: 10,
        });
        Object.defineProperty(event, 'target', { value: target });

        expect(wordAtPoint(event)).toBe('grace');
    });

    it('uses caretPositionFromPoint when available', () => {
        document.body.innerHTML =
            '<p data-verse="1"><span>the grace of God</span></p>';
        const textNode = document.querySelector('span')!.firstChild as Text;

        (
            document as Document & {
                caretPositionFromPoint: (
                    x: number,
                    y: number,
                ) => { offsetNode: Node; offset: number } | null;
            }
        ).caretPositionFromPoint = () => ({
            offsetNode: textNode,
            offset: 8,
        });

        const target = document.querySelector('span')!;
        const event = new MouseEvent('contextmenu', {
            bubbles: true,
            clientX: 10,
            clientY: 10,
        });
        Object.defineProperty(event, 'target', { value: target });

        expect(wordAtPoint(event)).toBe('grace');
    });
});
