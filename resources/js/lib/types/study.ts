import type { ViewMode } from '@/lib/types/bible';

export type StudySettings = {
    activeView: ViewMode;
    bookId: string;
    chapter: number;
    translationId: string;
    translationBId: string;
};
