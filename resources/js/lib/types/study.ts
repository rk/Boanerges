import type { ColumnContentType } from '@/lib/columns/catalog';

export type { ColumnContentType };

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
