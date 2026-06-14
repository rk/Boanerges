export type ChapterNavTarget = {
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
};

export type CatalogTranslation = Translation & {
    module: string;
    installed: boolean;
    bundled: boolean;
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
