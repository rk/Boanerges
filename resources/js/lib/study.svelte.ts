import { http } from '@inertiajs/svelte';

import { updateStudy as updateStudySettings } from '@/actions/App/Http/Controllers/SettingsController';
import { getAdjacentChapter } from '@/lib/bible.svelte.ts';
import type { StudySettings } from '@/lib/types/study';
import type { ViewMode } from '@/lib/types/bible';

export const study = $state({
    activeView: 'bible' as ViewMode,
    bookId: 'gen',
    chapter: 15,
    translationId: 'asv',
    translationBId: 'asv',
    scrollSync: false,
    settingsOpen: false,
});

let hydrated = false;
let persistTimeout: ReturnType<typeof setTimeout> | null = null;

export function hydrateStudy(settings: StudySettings): void {
    study.activeView = settings.activeView;
    study.bookId = settings.bookId;
    study.chapter = settings.chapter;
    study.translationId = settings.translationId;
    study.translationBId = settings.translationBId;
    hydrated = true;
}

function studyPayload(): StudySettings {
    return {
        activeView: study.activeView,
        bookId: study.bookId,
        chapter: study.chapter,
        translationId: study.translationId,
        translationBId: study.translationBId,
    };
}

function schedulePersist(): void {
    if (! hydrated) {
        return;
    }

    if (persistTimeout) {
        clearTimeout(persistTimeout);
    }

    persistTimeout = setTimeout(() => {
        http.patch(updateStudySettings.url(), studyPayload());
    }, 300);
}

export function getPreviousChapter(): ReturnType<typeof getAdjacentChapter> {
    return getAdjacentChapter(study.bookId, study.chapter, 'prev');
}

export function getNextChapter(): ReturnType<typeof getAdjacentChapter> {
    return getAdjacentChapter(study.bookId, study.chapter, 'next');
}

export function setView(view: ViewMode): void {
    study.activeView = view;
    schedulePersist();
}

export function setBook(nextBookId: string): void {
    study.bookId = nextBookId;
    study.chapter = 1;
    schedulePersist();
}

export function setChapter(nextChapter: number): void {
    study.chapter = nextChapter;
    schedulePersist();
}

export function setTranslation(id: string): void {
    study.translationId = id;
    schedulePersist();
}

export function setTranslationB(id: string): void {
    study.translationBId = id;
    schedulePersist();
}

export function goToPreviousChapter(): void {
    const previous = getPreviousChapter();

    if (previous) {
        study.bookId = previous.bookId;
        study.chapter = previous.chapter;
        schedulePersist();
    }
}

export function goToNextChapter(): void {
    const next = getNextChapter();

    if (next) {
        study.bookId = next.bookId;
        study.chapter = next.chapter;
        schedulePersist();
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

export function chapterLabel(
    book: string,
    chapter: number,
    translationAbbrev?: string,
): string {
    if (translationAbbrev) {
        return `${book} ${chapter} (${translationAbbrev})`;
    }

    return `${book} ${chapter}`;
}
