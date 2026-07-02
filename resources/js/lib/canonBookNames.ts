import type { Testament } from '@/lib/types/bible';

/** Canonical OSIS book IDs → English display names (matches OsisBookId::DISPLAY_NAMES). */
export const CANON_BOOK_NAMES: Record<string, string> = {
    gen: 'Genesis',
    exo: 'Exodus',
    lev: 'Leviticus',
    num: 'Numbers',
    deu: 'Deuteronomy',
    jos: 'Joshua',
    jdg: 'Judges',
    rut: 'Ruth',
    '1sa': '1 Samuel',
    '2sa': '2 Samuel',
    '1ki': '1 Kings',
    '2ki': '2 Kings',
    '1ch': '1 Chronicles',
    '2ch': '2 Chronicles',
    ezr: 'Ezra',
    neh: 'Nehemiah',
    est: 'Esther',
    job: 'Job',
    psa: 'Psalms',
    pro: 'Proverbs',
    ecc: 'Ecclesiastes',
    sng: 'Song of Songs',
    isa: 'Isaiah',
    jer: 'Jeremiah',
    lam: 'Lamentations',
    ezk: 'Ezekiel',
    dan: 'Daniel',
    hos: 'Hosea',
    jol: 'Joel',
    amo: 'Amos',
    oba: 'Obadiah',
    jon: 'Jonah',
    mic: 'Micah',
    nam: 'Nahum',
    hab: 'Habakkuk',
    zep: 'Zephaniah',
    hag: 'Haggai',
    zec: 'Zechariah',
    mal: 'Malachi',
    mat: 'Matthew',
    mrk: 'Mark',
    luk: 'Luke',
    jhn: 'John',
    act: 'Acts',
    rom: 'Romans',
    '1co': '1 Corinthians',
    '2co': '2 Corinthians',
    gal: 'Galatians',
    eph: 'Ephesians',
    php: 'Philippians',
    col: 'Colossians',
    '1th': '1 Thessalonians',
    '2th': '2 Thessalonians',
    '1ti': '1 Timothy',
    '2ti': '2 Timothy',
    tit: 'Titus',
    phm: 'Philemon',
    heb: 'Hebrews',
    jas: 'James',
    '1pe': '1 Peter',
    '2pe': '2 Peter',
    '1jn': '1 John',
    '2jn': '2 John',
    '3jn': '3 John',
    jud: 'Jude',
    rev: 'Revelation',
};

const CANON_BOOK_IDS = Object.keys(CANON_BOOK_NAMES);
const NT_START_INDEX = CANON_BOOK_IDS.indexOf('mat');

export const CANON_OT_BOOK_IDS: readonly string[] = CANON_BOOK_IDS.slice(
    0,
    NT_START_INDEX,
);
export const CANON_NT_BOOK_IDS: readonly string[] =
    CANON_BOOK_IDS.slice(NT_START_INDEX);

export function canonBookName(bookId: string): string {
    return CANON_BOOK_NAMES[bookId.toLowerCase()] ?? bookId.toUpperCase();
}

export function canonBookTestament(bookId: string): Testament {
    const id = bookId.toLowerCase();
    const index = CANON_BOOK_IDS.indexOf(id);

    return index >= NT_START_INDEX ? 'nt' : 'ot';
}
