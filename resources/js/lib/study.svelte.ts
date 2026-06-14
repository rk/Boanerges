import { getAdjacentChapter, getChapter } from '@/lib/mock/chapter';
import type { Chapter, ViewMode } from '@/lib/types/bible';

export const study = $state({
    activeView: 'bible' as ViewMode,
    bookId: 'gen',
    chapter: 15,
    translationId: 'kjv',
    translationBId: 'asv',
    scrollSync: false,
    settingsOpen: false,
});

export function getCurrentChapter(): Chapter {
    return getChapter(study.bookId, study.chapter);
}

export function getPreviousChapter(): ReturnType<typeof getAdjacentChapter> {
    return getAdjacentChapter(study.bookId, study.chapter, 'prev');
}

export function getNextChapter(): ReturnType<typeof getAdjacentChapter> {
    return getAdjacentChapter(study.bookId, study.chapter, 'next');
}

export function setView(view: ViewMode): void {
    study.activeView = view;
}

export function setBook(nextBookId: string): void {
    study.bookId = nextBookId;
    study.chapter = 1;
}

export function setChapter(nextChapter: number): void {
    study.chapter = nextChapter;
}

export function setTranslation(id: string): void {
    study.translationId = id;
}

export function setTranslationB(id: string): void {
    study.translationBId = id;
}

export function goToPreviousChapter(): void {
    const previous = getPreviousChapter();

    if (previous) {
        study.bookId = previous.bookId;
        study.chapter = previous.chapter;
    }
}

export function goToNextChapter(): void {
    const next = getNextChapter();

    if (next) {
        study.bookId = next.bookId;
        study.chapter = next.chapter;
    }
}

export function openSettings(): void {
    study.settingsOpen = true;
}

export function closeSettings(): void {
    study.settingsOpen = false;
}

export function setScrollSync(enabled: boolean): void {
    study.scrollSync = enabled;
}

export function navLabel(bookAbbrev: string, chapterNumber: number, direction: 'prev' | 'next'): string {
    const arrow = direction === 'prev' ? '^' : 'v';

    return `${arrow} ${bookAbbrev} ${chapterNumber} ${arrow}`;
}

export function chapterLabel(current: Chapter, translationAbbrev?: string): string {
    if (translationAbbrev) {
        return `${current.book} ${current.chapter} (${translationAbbrev})`;
    }

    return `${current.book} ${current.chapter}`;
}
