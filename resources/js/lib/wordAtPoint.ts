// ponytail: debug order tries caretRangeFromPoint first; check console for
// `[wordAtPoint] { method, word }` then delete unused fallbacks.
const WORD_CHAR_PATTERN = /[\p{L}\p{N}'-]+/u;
const WORD_END_PATTERN = /^[\p{L}\p{N}'-]+/u;
const WORD_START_PATTERN = /[\p{L}\p{N}'-]+$/u;

export type WordLookupMethod =
    | 'text-selection'
    | 'caretRangeFromPoint'
    | 'caretPositionFromPoint'
    | 'elementFromPoint'
    | 'data-dict-word'
    | 'none';

type DocumentWithCaret = Document & {
    caretRangeFromPoint?: (x: number, y: number) => Range | null;
    caretPositionFromPoint?: (
        x: number,
        y: number,
    ) => { offsetNode: Node; offset: number } | null;
};

let lastWordLookupMethod: WordLookupMethod = 'none';

export function getLastWordLookupMethod(): WordLookupMethod {
    return lastWordLookupMethod;
}

export function logWordLookup(method: WordLookupMethod, word: string): void {
    lastWordLookupMethod = method;
    console.debug('[wordAtPoint]', { method, word });
}

export function extractWordFromTextOffset(
    text: string,
    offset: number,
): string {
    const before = text.slice(0, offset);
    const after = text.slice(offset);
    const start = before.search(WORD_START_PATTERN);
    const endMatch = after.match(WORD_END_PATTERN);
    const end = endMatch ? offset + endMatch[0].length : offset;

    if (start === -1) {
        return '';
    }

    return text.slice(start, end);
}

function rangeFromCaretRangeFromPoint(x: number, y: number): Range | null {
    return (document as DocumentWithCaret).caretRangeFromPoint?.(x, y) ?? null;
}

function rangeFromCaretPositionFromPoint(x: number, y: number): Range | null {
    const position = (document as DocumentWithCaret).caretPositionFromPoint?.(
        x,
        y,
    );

    if (!position) {
        return null;
    }

    const range = document.createRange();
    range.setStart(position.offsetNode, position.offset);
    range.collapse(true);

    return range;
}

function pointInRect(x: number, y: number, rect: DOMRect): boolean {
    return (
        x >= rect.left && x <= rect.right && y >= rect.top && y <= rect.bottom
    );
}

function textOffsetAtCoordinates(textNode: Text, x: number, y: number): number {
    const text = textNode.data;

    if (text.length === 0) {
        return 0;
    }

    const range = document.createRange();
    let left = 0;
    let right = text.length;

    while (left < right) {
        const mid = Math.floor((left + right) / 2);
        range.setStart(textNode, 0);
        range.setEnd(textNode, mid);
        const rect = range.getBoundingClientRect();

        if (x > rect.right || (x === rect.right && y > rect.bottom)) {
            left = mid + 1;
        } else {
            right = mid;
        }
    }

    return left;
}

function findTextNodeAtCoordinates(
    root: Element,
    x: number,
    y: number,
): { node: Text; offset: number } | null {
    const walker = document.createTreeWalker(root, NodeFilter.SHOW_TEXT);
    let textNode = walker.nextNode() as Text | null;

    while (textNode) {
        if (textNode.data !== '') {
            const range = document.createRange();
            range.selectNodeContents(textNode);
            const rects = range.getClientRects();

            for (let index = 0; index < rects.length; index++) {
                if (pointInRect(x, y, rects[index])) {
                    return {
                        node: textNode,
                        offset: textOffsetAtCoordinates(textNode, x, y),
                    };
                }
            }
        }

        textNode = walker.nextNode() as Text | null;
    }

    return null;
}

function rangeFromElementFromPoint(
    x: number,
    y: number,
    root: Element | null,
): Range | null {
    const element = document.elementFromPoint(x, y);

    if (!element || (root && !root.contains(element))) {
        return null;
    }

    const searchRoot = root ?? element.closest('[data-verse]') ?? element;
    const located = findTextNodeAtCoordinates(searchRoot, x, y);

    if (!located) {
        return null;
    }

    const range = document.createRange();
    range.setStart(located.node, located.offset);
    range.collapse(true);

    return range;
}

function wordFromDictWordElement(target: EventTarget | null): string {
    const element = target instanceof Element ? target : null;
    const wordElement = element?.closest('[data-dict-word]');

    if (!(wordElement instanceof HTMLElement)) {
        return '';
    }

    return wordElement.dataset.dictWord ?? wordElement.textContent ?? '';
}

function wordFromRange(range: Range, method: WordLookupMethod): string {
    if (!range.startContainer.textContent) {
        return '';
    }

    return extractWordFromTextOffset(
        range.startContainer.textContent,
        range.startOffset,
    );
}

export function wordAtPoint(event: MouseEvent): string {
    const root =
        event.target instanceof Element
            ? event.target.closest('[data-verse]')
            : null;
    const x = event.clientX;
    const y = event.clientY;

    const caretRange = rangeFromCaretRangeFromPoint(x, y);

    if (caretRange) {
        const word = wordFromRange(caretRange, 'caretRangeFromPoint');

        if (word !== '') {
            logWordLookup('caretRangeFromPoint', word);

            return word;
        }
    }

    const caretPositionRange = rangeFromCaretPositionFromPoint(x, y);

    if (caretPositionRange) {
        const word = wordFromRange(
            caretPositionRange,
            'caretPositionFromPoint',
        );

        if (word !== '') {
            logWordLookup('caretPositionFromPoint', word);

            return word;
        }
    }

    const elementRange = rangeFromElementFromPoint(x, y, root);

    if (elementRange) {
        const word = wordFromRange(elementRange, 'elementFromPoint');

        if (word !== '') {
            logWordLookup('elementFromPoint', word);

            return word;
        }
    }

    const dictWord = wordFromDictWordElement(event.target);

    if (dictWord !== '') {
        const word = dictWord.match(WORD_CHAR_PATTERN)?.[0] ?? '';

        if (word !== '') {
            logWordLookup('data-dict-word', word);

            return word;
        }
    }

    logWordLookup('none', '');

    return '';
}
