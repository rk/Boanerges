export type VerseHighlight = {
    verse: number;
    endVerse: number | null;
};

export function verseNumbersInRange(verse: number, endVerse: number | null): Set<number> {
    const last = endVerse ?? verse;
    const start = Math.min(verse, last);
    const end = Math.max(verse, last);
    const numbers = new Set<number>();

    for (let number = start; number <= end; number++) {
        numbers.add(number);
    }

    return numbers;
}
