<script lang="ts">
    import Book from '@lucide/svelte/icons/book';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import ParagraphText from '@/components/reader/ParagraphText.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import VerseText from '@/components/reader/VerseText.svelte';
    import TranslationSelect from '@/components/sidebar/TranslationSelect.svelte';
    import { bible, bookAbbrev, fetchChapter, getAdjacentChapter, peekChapter } from '@/lib/bible.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import {
        goToNextChapter,
        goToPreviousChapter,
        study,
    } from '@/lib/study.svelte.ts';
    import type { Chapter, ChapterNavTarget } from '@/lib/types/bible';
    import { showVerseContextMenu } from '@/lib/verseContextMenu';
    import { verseNumbersInRange } from '@/lib/verseHighlight';

    let {
        translationId,
        variant = 'paragraph',
        showChapterNav = false,
        slotIndex,
        onTranslationChange,
        scrollRef = $bindable<HTMLElement | null>(null),
        onscroll,
    }: {
        translationId: string;
        variant?: 'paragraph' | 'verse';
        showChapterNav?: boolean;
        slotIndex?: number;
        onTranslationChange?: (id: string) => void;
        scrollRef?: HTMLElement | null;
        onscroll?: () => void;
    } = $props();

    let currentChapter = $state<Chapter | null>(null);
    let loading = $state(true);
    let lastLocationKey: string | null = null;
    let lastTranslationId: string | null = null;

    const readerStyle = $derived(getReaderStyle());
    const translation = $derived(bible.translations.find((item) => item.id === translationId));
    const highlightedVerses = $derived(
        study.verseHighlight
            ? verseNumbersInRange(study.verseHighlight.verse, study.verseHighlight.endVerse)
            : new Set<number>(),
    );
    const isPrimary = $derived(slotIndex === undefined);
    const contentType = $derived(isPrimary ? 'primary-bible' as const : 'bible-secondary' as const);

    const previousChapter = $derived(getAdjacentChapter(study.bookId, study.chapter, 'prev'));
    const nextChapter = $derived(getAdjacentChapter(study.bookId, study.chapter, 'next'));

    $effect(() => {
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        const locationKey = `${bookId}:${chapterNumber}`;
        const navigated = lastLocationKey !== null && locationKey !== lastLocationKey;
        const translationChanged = lastTranslationId !== null && translationId !== lastTranslationId;
        const isInitialLoad = lastLocationKey === null;
        let cancelled = false;

        lastLocationKey = locationKey;
        lastTranslationId = translationId;

        if (isInitialLoad || navigated) {
            loading = true;
            currentChapter = null;
        } else if (translationChanged) {
            const cached = peekChapter(translationId, bookId, chapterNumber);

            if (cached) {
                currentChapter = cached;
                loading = false;
            }
        }

        void fetchChapter(translationId, bookId, chapterNumber)
            .then((chapter) => {
                if (! cancelled) {
                    currentChapter = chapter;
                    loading = false;
                }
            })
            .catch(() => {
                if (! cancelled) {
                    loading = false;
                }
            });

        return () => {
            cancelled = true;
        };
    });

    function chapterNavTarget(bookId: string, chapterNumber: number): ChapterNavTarget {
        return {
            bookId,
            bookAbbrev: bookAbbrev(bookId),
            chapter: chapterNumber,
        };
    }

    const prevNav = $derived(
        showChapterNav && previousChapter
            ? chapterNavTarget(previousChapter.bookId, previousChapter.chapter)
            : undefined,
    );

    const nextNav = $derived(
        showChapterNav && nextChapter
            ? chapterNavTarget(nextChapter.bookId, nextChapter.chapter)
            : undefined,
    );

    function handleVerseContextMenu(event: MouseEvent): void {
        const target = event.target as HTMLElement;
        const verseEl = target.closest('[data-verse]');

        if (! verseEl) {
            return;
        }

        const verse = Number(verseEl.getAttribute('data-verse'));

        if (! Number.isFinite(verse)) {
            return;
        }

        showVerseContextMenu(verse, event);
    }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" style={readerStyle}>
    {#if onTranslationChange}
        <ColumnHeader
            {contentType}
            {slotIndex}
            showViewSelector={! isPrimary}
        >
            <Book size={14} class="text-base-content/70 shrink-0" aria-hidden="true" />
            <TranslationSelect
                value={translationId}
                onchange={onTranslationChange}
            />
        </ColumnHeader>
    {/if}

    {#if loading || ! currentChapter}
        <div class="flex flex-1 items-center justify-center p-8">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>
    {:else}
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <div class="flex min-h-0 flex-1 flex-col" oncontextmenu={handleVerseContextMenu}>
            <ReaderPane
                chapter={currentChapter}
                translationAbbrev={translation?.abbrev}
                {prevNav}
                {nextNav}
                onprev={prevNav ? goToPreviousChapter : undefined}
                onnext={nextNav ? goToNextChapter : undefined}
                bind:scrollRef
                {onscroll}
            >
                {#if variant === 'verse'}
                    <VerseText verses={currentChapter.verses} {highlightedVerses} />
                {:else}
                    <ParagraphText verses={currentChapter.verses} {highlightedVerses} />
                {/if}
            </ReaderPane>
        </div>
    {/if}
</div>
