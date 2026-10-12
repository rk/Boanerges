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
    const originalCaretRangeFromPoint = (
        document as Document & {
            caretRangeFromPoint?: (x: number, y: number) => Range | null;
        }
    ).caretRangeFromPoint?.bind(document);

    afterEach(() => {
        document.body.innerHTML = '';

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

    it('reads the word at the caret range', () => {
        document.body.innerHTML =
            '<p data-verse="1"><span>the grace of God</span></p>';
        const textNode = document.querySelector('span')!.firstChild as Text;
        const target = document.querySelector('span')!;
        const range = document.createRange();
        range.setStart(textNode, 8);
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

        expect(wordAtPoint(event)).toBe('grace');
    });
});
