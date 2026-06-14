<script lang="ts">
    import { fetchChapter, getAdjacentChapter, loadBooks, bookAbbrev } from '@/lib/bible.svelte.ts';
    import {
        goToNextChapter,
        goToPreviousChapter,
        study,
    } from '@/lib/study.svelte.ts';
    import ParagraphText from '@/components/reader/ParagraphText.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import type { Chapter, ChapterNavTarget } from '@/lib/types/bible';

    let currentChapter = $state<Chapter | null>(null);
    let loading = $state(true);

    const previousChapter = $derived(getAdjacentChapter(study.bookId, study.chapter, 'prev'));
    const nextChapter = $derived(getAdjacentChapter(study.bookId, study.chapter, 'next'));
    const readerStyle = $derived(getReaderStyle());

    $effect(() => {
        const translationId = study.translationId;
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        let cancelled = false;

        loading = true;
        currentChapter = null;

        loadBooks(translationId)
            .then(() => fetchChapter(translationId, bookId, chapterNumber))
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
            bookAbbrev: bookAbbrev(bookId),
            chapter: chapterNumber,
        };
    }

    const prevNav = $derived(
        previousChapter
            ? chapterNavTarget(previousChapter.bookId, previousChapter.chapter)
            : undefined,
    );

    const nextNav = $derived(
        nextChapter ? chapterNavTarget(nextChapter.bookId, nextChapter.chapter) : undefined,
    );
</script>

<div class="h-full min-h-0" style={readerStyle}>
    <div class="mx-auto h-full max-w-prose">
        {#if loading || ! currentChapter}
            <div class="flex h-full items-center justify-center p-8">
                <span class="loading loading-spinner loading-lg text-primary"></span>
            </div>
        {:else}
            <ReaderPane
                chapter={currentChapter}
                {prevNav}
                {nextNav}
                onprev={previousChapter ? goToPreviousChapter : undefined}
                onnext={nextChapter ? goToNextChapter : undefined}
            >
                <ParagraphText verses={currentChapter.verses} />
            </ReaderPane>
        {/if}
    </div>
</div>
