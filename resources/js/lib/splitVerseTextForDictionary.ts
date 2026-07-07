export type VerseTextPart = {
    kind: 'word' | 'text';
    value: string;
};

const VERSE_TEXT_PART_PATTERN = /[\p{L}\p{N}'-]+|[^\p{L}\p{N}'-]+/gu;
const WORD_CHAR_PATTERN = /[\p{L}\p{N}'-]/u;

export function splitVerseTextForDictionary(text: string): VerseTextPart[] {
    if (text === '') {
        return [];
    }

    const parts: VerseTextPart[] = [];

    for (const value of text.match(VERSE_TEXT_PART_PATTERN) ?? []) {
        parts.push({
            kind: WORD_CHAR_PATTERN.test(value) ? 'word' : 'text',
            value,
        });
    }

    return parts;
}
