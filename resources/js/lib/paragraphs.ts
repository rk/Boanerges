import type { Verse } from '@/lib/types/bible';

export function groupVersesIntoParagraphs(verses: Verse[]): Verse[][] {
    const groups: Verse[][] = [];
    let current: Verse[] = [];

    for (const verse of verses) {
        if (verse.paragraphStart && current.length > 0) {
            groups.push(current);
            current = [];
        }

        current.push(verse);
    }

    if (current.length > 0) {
        groups.push(current);
    }

    return groups;
}

export function splitTextParagraphs(text: string): string[] {
    if (text.trim() === '') {
        return [];
    }

    return text
        .split(/\n\s*\n/)
        .map((part) => part.trim())
        .filter((part) => part.length > 0);
}

export type ScribePreviewSegment = {
    verseNumber: number;
    text: string;
    showVerseNumber: boolean;
};

export function buildScribePreviewParagraphs(
    verses: Verse[],
): ScribePreviewSegment[][] {
    const groups = groupVersesIntoParagraphs(verses);
    const result: ScribePreviewSegment[][] = [];

    for (const group of groups) {
        let currentParagraph: ScribePreviewSegment[] = [];

        for (const verse of group) {
            const parts = splitTextParagraphs(verse.text);

            for (let index = 0; index < parts.length; index++) {
                if (index > 0) {
                    if (currentParagraph.length > 0) {
                        result.push(currentParagraph);
                        currentParagraph = [];
                    }
                }

                currentParagraph.push({
                    verseNumber: verse.number,
                    text: parts[index],
                    showVerseNumber: index === 0,
                });
            }
        }

        if (currentParagraph.length > 0) {
            result.push(currentParagraph);
        }
    }

    return result;
}

export function toPreviewVerses(
    sourceVerses: Verse[],
    getText: (verseNumber: number) => string,
    getParagraphStart: (verseNumber: number) => boolean,
): Verse[] {
    const result: Verse[] = [];
    let pendingBreak = false;

    for (const source of sourceVerses) {
        const text = getText(source.number);
        const hasText = text.trim().length > 0;
        const startsParagraph =
            pendingBreak || getParagraphStart(source.number);

        pendingBreak = false;

        if (!hasText) {
            if (startsParagraph && result.length > 0) {
                pendingBreak = true;
            }

            continue;
        }

        result.push({
            number: source.number,
            text,
            paragraphStart: result.length === 0 ? true : startsParagraph,
        });
    }

    return result;
}
