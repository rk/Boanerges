import { SvelteSet } from 'svelte/reactivity';
import { updateStudy as updateStudySettings } from '@/actions/App/Http/Controllers/SettingsController';
import { getAdjacentChapter, bible } from '@/lib/bible.svelte.ts';
import { setComparisonInput } from '@/lib/comparison.svelte.ts';
import { setCrossReferenceInput } from '@/lib/crossrefs.svelte.ts';
import { setDictionaryWord } from '@/lib/dictionary.svelte.ts';
import { patchJson } from '@/lib/patchJson';
import { formatScriptureReference } from '@/lib/scriptureReference';
import type { ScriptureReference } from '@/lib/scriptureReference';
import {
    comparisonTargetSlot,
    crossReferencesTargetSlot,
    dictionaryTargetSlot,
    normalizeColumns,
    sanitizeStudySettings,
    verseListTargetSlot,
} from '@/lib/studyLayout';
import type { ColumnContentType, StudySettings } from '@/lib/types/study';
import type { VerseHighlight } from '@/lib/verseHighlight';
import {
    addVerseToList,
    hydrateVerseListSettings,
} from '@/lib/verseList.svelte.ts';

export const study = $state({
    columnCount: 1 as 1 | 2 | 3,
    columns: [] as ColumnContentType[],
    bookId: 'gen',
    chapter: 15,
    translationId: 'asv',
    translationBId: 'asv',
    translationCId: 'asv',
    verseListShowContent: true,
    verseListActiveId: null as string | null,
    scrollSync: false,
    settingsOpen: false,
    verseHighlight: null as VerseHighlight | null,
});

let hydrated = false;
let persistTimeout: ReturnType<typeof setTimeout> | null = null;

export function hydrateStudy(settings: StudySettings): void {
    if (hydrated) {
        return;
    }

    const sanitized = sanitizeStudySettings(settings);

    study.columnCount = sanitized.columnCount;
    study.columns = sanitized.columns;
    study.bookId = sanitized.bookId;
    study.chapter = sanitized.chapter;
    study.translationId = sanitized.translationId;
    study.translationBId = sanitized.translationBId;
    study.translationCId = sanitized.translationCId;
    study.verseListShowContent = sanitized.verseListShowContent ?? true;
    study.verseListActiveId = sanitized.verseListActiveId ?? null;
    hydrateVerseListSettings(
        study.verseListShowContent,
        study.verseListActiveId,
    );
    hydrated = true;
}

function studyPayload(): StudySettings {
    return {
        columnCount: study.columnCount,
        columns: normalizeColumns(study.columnCount, study.columns),
        bookId: study.bookId,
        chapter: study.chapter,
        translationId: study.translationId,
        translationBId: study.translationBId,
        translationCId: study.translationCId,
        verseListShowContent: study.verseListShowContent,
        verseListActiveId: study.verseListActiveId,
    };
}

function schedulePersist(): void {
    if (!hydrated) {
        return;
    }

    if (persistTimeout) {
        clearTimeout(persistTimeout);
    }

    persistTimeout = setTimeout(() => {
        void patchJson(updateStudySettings.url(), studyPayload());
    }, 300);
}

export function getPreviousChapter(): ReturnType<typeof getAdjacentChapter> {
    return getAdjacentChapter(study.bookId, study.chapter, 'prev');
}

export function getNextChapter(): ReturnType<typeof getAdjacentChapter> {
    return getAdjacentChapter(study.bookId, study.chapter, 'next');
}

export function setColumnCount(count: 1 | 2 | 3): void {
    study.columnCount = count;
    study.columns = normalizeColumns(count, study.columns);
    schedulePersist();
}

export function setColumnContent(
    slotIndex: number,
    type: ColumnContentType,
): void {
    const columns = normalizeColumns(study.columnCount, [...study.columns]);
    columns[slotIndex] = type;
    study.columns = columns;
    schedulePersist();
}

export function setBook(nextBookId: string): void {
    study.bookId = nextBookId;
    study.chapter = 1;
    study.verseHighlight = null;
    schedulePersist();
}

export function setChapter(nextChapter: number): void {
    study.chapter = nextChapter;
    study.verseHighlight = null;
    schedulePersist();
}

export function goToVerseReference(
    bookId: string,
    chapter: number,
    verse: number,
    endVerse: number | null = null,
): void {
    study.bookId = bookId;
    study.chapter = chapter;
    study.verseHighlight = { verse, endVerse };
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

export function setTranslationC(id: string): void {
    study.translationCId = id;
    schedulePersist();
}

export function goToPreviousChapter(): void {
    const previous = getPreviousChapter();

    if (previous) {
        study.bookId = previous.bookId;
        study.chapter = previous.chapter;
        study.verseHighlight = null;
        schedulePersist();
    }
}

export function goToNextChapter(): void {
    const next = getNextChapter();

    if (next) {
        study.bookId = next.bookId;
        study.chapter = next.chapter;
        study.verseHighlight = null;
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

export function syncStudyTranslationSelection(): void {
    const installedIds = new SvelteSet(
        bible.translations.map((translation) => translation.id),
    );

    if (!installedIds.has(study.translationId)) {
        study.translationId = 'asv';
    }

    if (!installedIds.has(study.translationBId)) {
        study.translationBId = study.translationId;
    }

    if (!installedIds.has(study.translationCId)) {
        study.translationCId = study.translationBId;
    }

    schedulePersist();
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

export function ensureSearchColumn(): void {
    if (study.columns.includes('search')) {
        return;
    }

    if (study.columnCount === 1) {
        setColumnCount(2);
    }

    const slotIndex = study.columns.findIndex(
        (column) => column !== 'bible-secondary',
    );

    if (slotIndex >= 0) {
        setColumnContent(slotIndex, 'search');

        return;
    }

    setColumnContent(study.columns.length - 1, 'search');
}

export function ensureCrossReferencesColumn(reference?: string): void {
    if (study.columnCount === 1) {
        setColumnCount(2);
    }

    const slotIndex = crossReferencesTargetSlot(
        study.columnCount,
        study.columns,
    );

    if (study.columns[slotIndex] !== 'cross-references') {
        setColumnContent(slotIndex, 'cross-references');
    }

    const resolvedReference =
        reference ??
        formatScriptureReference(
            study.bookId,
            study.chapter,
            study.verseHighlight?.verse ?? 1,
            bible.books,
        );

    setCrossReferenceInput(resolvedReference);
}

export function ensureDictionaryColumn(word?: string): void {
    const existingSlot = study.columns.findIndex(
        (column) => column === 'dictionary',
    );

    if (existingSlot >= 0) {
        if (word !== undefined && word.trim() !== '') {
            setDictionaryWord(word);
        }

        return;
    }

    if (study.columnCount === 1) {
        setColumnCount(2);
    }

    const slotIndex = dictionaryTargetSlot(study.columnCount, study.columns);

    if (study.columns[slotIndex] !== 'dictionary') {
        setColumnContent(slotIndex, 'dictionary');
    }

    if (word !== undefined && word.trim() !== '') {
        setDictionaryWord(word);
    }
}

export function setVerseListShowContentSetting(enabled: boolean): void {
    study.verseListShowContent = enabled;
    schedulePersist();
}

export function setVerseListActiveId(id: string | null): void {
    study.verseListActiveId = id;
    schedulePersist();
}

export function ensureComparisonColumn(reference?: string): void {
    if (study.columnCount === 1) {
        setColumnCount(2);
    }

    const slotIndex = comparisonTargetSlot(study.columnCount, study.columns);

    if (study.columns[slotIndex] !== 'comparison') {
        setColumnContent(slotIndex, 'comparison');
    }

    const resolvedReference =
        reference ??
        formatScriptureReference(
            study.bookId,
            study.chapter,
            study.verseHighlight?.verse ?? 1,
            bible.books,
        );

    setComparisonInput(resolvedReference);
}

export function ensureVerseListColumn(ref?: ScriptureReference): void {
    if (study.columnCount === 1) {
        setColumnCount(2);
    }

    const slotIndex = verseListTargetSlot(study.columnCount, study.columns);

    if (study.columns[slotIndex] !== 'verse-list') {
        setColumnContent(slotIndex, 'verse-list');
    }

    if (ref) {
        addVerseToList(ref, bible.books);
    }
}
