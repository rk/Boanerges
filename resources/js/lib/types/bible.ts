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
};

export type ChapterNav = {
    book: string;
    bookAbbrev: string;
    chapter: number;
};
