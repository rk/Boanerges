export type ColumnContentType = 'bible-secondary' | 'notes' | 'scribe' | 'search' | 'cross-references';

export type StudySettings = {
    columnCount: 1 | 2 | 3;
    columns: ColumnContentType[];
    bookId: string;
    chapter: number;
    translationId: string;
    translationBId: string;
    translationCId: string;
};
