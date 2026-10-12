import type { Verse } from '@/lib/types/bible';

export type ScribeDraftEntry = {
    text: string;
    paragraphStartOverride?: boolean;
};

export function effectiveParagraphStart(
    verseNumber: number,
    override?: boolean,
): boolean {
    if (override !== undefined) {
        return override;
    }

    return verseNumber === 1;
}

export function applyPrimaryParagraphBreaks(
    verseNumbers: number[],
    entries: Record<number, ScribeDraftEntry>,
    sourceVerses: Verse[],
): Record<number, ScribeDraftEntry> {
    const updated = { ...entries };

    for (const verseNumber of verseNumbers) {
        if (verseNumber === 1) {
            continue;
        }

        const source = sourceVerses.find(
            (verse) => verse.number === verseNumber,
        );
        const fromPrimary = source?.paragraphStart ?? false;
        const entry = updated[verseNumber] ?? { text: '' };

        if (fromPrimary) {
            updated[verseNumber] = {
                ...entry,
                paragraphStartOverride: true,
            };
        } else {
            const { paragraphStartOverride: _, ...rest } = entry;

            if (rest.text === undefined || rest.text === '') {
                delete updated[verseNumber];
            } else {
                updated[verseNumber] = rest;
            }
        }
    }

    return updated;
}

export function resetInterVerseParagraphBreaks(
    verseNumbers: number[],
    entries: Record<number, ScribeDraftEntry>,
): Record<number, ScribeDraftEntry> {
    const updated = { ...entries };

    for (const verseNumber of verseNumbers) {
        if (verseNumber === 1) {
            continue;
        }

        const entry = updated[verseNumber];

        if (entry?.paragraphStartOverride === undefined) {
            continue;
        }

        const { paragraphStartOverride: _, ...rest } = entry;

        if (rest.text === undefined || rest.text === '') {
            delete updated[verseNumber];
        } else {
            updated[verseNumber] = rest;
        }
    }

    return updated;
}
