import type { Book, CrossReference } from '@/lib/types/bible';
import { canonBookName } from '@/lib/canonBookNames';

export type ScriptureReference = {
    bookId: string;
    chapter: number;
    verse: number;
};

export function parseScriptureReference(input: string, books: Book[]): ScriptureReference | null {
    const trimmed = input.trim();

    if (trimmed === '') {
        return null;
    }

    const match = trimmed.match(/^(.+)\s+(\d+)\s*:\s*(\d+)\s*$/);

    if (match === null) {
        return null;
    }

    const bookPart = normalizeReferenceToken(match[1]);
    const chapter = Number.parseInt(match[2], 10);
    const verse = Number.parseInt(match[3], 10);

    if (chapter <= 0 || verse <= 0) {
        return null;
    }

    const book = resolveBook(bookPart, books);

    if (book === null) {
        return null;
    }

    return { bookId: book.id, chapter, verse };
}

export function formatScriptureReference(
    bookId: string,
    chapter: number,
    verse: number,
    books: Book[],
): string {
    const name = bookDisplayName(bookId, books);

    return `${name} ${chapter}:${verse}`;
}

export function formatCrossReference(ref: CrossReference, books: Book[]): string {
    const name = ref.bookName ?? bookDisplayName(ref.bookId, books);
    const range = ref.endVerse !== null && ref.endVerse !== ref.verse
        ? `${ref.verse}–${ref.endVerse}`
        : `${ref.verse}`;

    return `${name} ${ref.chapter}:${range}`;
}

function bookDisplayName(bookId: string, books: Book[]): string {
    return books.find((entry) => entry.id === bookId)?.name ?? canonBookName(bookId);
}

function normalizeReferenceToken(value: string): string {
    return value.trim().toLowerCase().replace(/\./g, '').replace(/\s+/g, ' ');
}

function resolveBook(bookPart: string, books: Book[]): Book | null {
    const normalized = normalizeReferenceToken(bookPart);

    const candidates = [...books].sort(
        (left, right) => bookMatchKeys(right).length - bookMatchKeys(left).length,
    );

    for (const book of candidates) {
        if (bookMatchKeys(book).some((key) => key === normalized)) {
            return book;
        }
    }

    return null;
}

function bookMatchKeys(book: Book): string[] {
    return [
        normalizeReferenceToken(book.name),
        normalizeReferenceToken(book.id),
        normalizeReferenceToken(book.abbrev),
        normalizeReferenceToken(canonBookName(book.id)),
    ];
}
