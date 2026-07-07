const WORD_END_PATTERN = /^[\p{L}\p{N}'-]+/u;
const WORD_START_PATTERN = /[\p{L}\p{N}'-]+$/u;

type DocumentWithCaret = Document & {
    caretRangeFromPoint?: (x: number, y: number) => Range | null;
};

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

export function wordAtPoint(event: MouseEvent): string {
    const range = (document as DocumentWithCaret).caretRangeFromPoint?.(
        event.clientX,
        event.clientY,
    );

    if (!range || !range.startContainer.textContent) {
        return '';
    }

    return extractWordFromTextOffset(
        range.startContainer.textContent,
        range.startOffset,
    );
}
