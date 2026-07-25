export type ColumnContentType =
    | 'bible-secondary'
    | 'notes'
    | 'scribe'
    | 'search'
    | 'cross-references'
    | 'dictionary'
    | 'comparison'
    | 'verse-list';

export type StudySettings = {
    columnCount: 1 | 2 | 3;
    columns: ColumnContentType[];
    bookId: string;
    chapter: number;
    translationId: string;
    translationBId: string;
    translationCId: string;
    verseListShowContent?: boolean;
    verseListActiveId?: string | null;
};
