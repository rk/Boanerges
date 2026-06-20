export type ChapterNavTarget = {
    bookId: string;
    bookAbbrev: string;
    chapter: number;
};

export type ViewMode = 'bible' | 'comparison' | 'scribe';

export type Testament = 'ot' | 'nt';

export type Verse = {
    number: number;
    text: string;
    paragraphStart?: boolean;
};

export type Chapter = {
    book: string;
    bookAbbrev: string;
    chapter: number;
    verses: Verse[];
};

export type Book = {
    id: string;
    name: string;
    abbrev: string;
    testament: Testament;
    chapters: number;
};

export type Translation = {
    id: string;
    name: string;
    abbrev: string;
    bundled?: boolean;
    install_status?: string | null;
    install_step?: string | null;
    install_error?: string | null;
};

export type CatalogTranslation = Translation & {
    module: string;
    installed: boolean;
    bundled: boolean;
    about: string;
};

export type SearchResult = {
    bookId: string;
    chapter: number;
    verse: number;
    snippet: string;
    translation: string;
};

export type CrossReference = {
    rank: number;
    bookId: string;
    bookName?: string;
    chapter: number;
    verse: number;
    endVerse: number | null;
};

export type ChapterNav = {
    book: string;
    bookAbbrev: string;
    chapter: number;
};

export type ScribeVerse = {
    verse: number;
    text: string;
    paragraphStart?: boolean;
};
