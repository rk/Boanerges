<script lang="ts">
    import ParagraphText from '@/components/reader/ParagraphText.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import { getChapter } from '@/lib/mock/chapter';
    import {
        getCurrentChapter,
        getNextChapter,
        getPreviousChapter,
        goToNextChapter,
        goToPreviousChapter,
    } from '@/lib/study.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import type { ChapterNavTarget } from '@/lib/types/bible';

    const currentChapter = $derived(getCurrentChapter());
    const previousChapter = $derived(getPreviousChapter());
    const nextChapter = $derived(getNextChapter());
    const readerStyle = $derived(getReaderStyle());

    function chapterNavTarget(
        bookId: string,
        chapterNumber: number,
    ): ChapterNavTarget {
        const chapter = getChapter(bookId, chapterNumber);

        return {
            bookAbbrev: chapter.bookAbbrev,
            chapter: chapterNumber,
        };
    }

    const prevNav = $derived(
        previousChapter
            ? chapterNavTarget(previousChapter.bookId, previousChapter.chapter)
            : undefined,
    );

    const nextNav = $derived(
        nextChapter
            ? chapterNavTarget(nextChapter.bookId, nextChapter.chapter)
            : undefined,
    );
</script>

<div class="h-full min-h-0" style={readerStyle}>
    <div class="mx-auto h-full max-w-prose">
        <ReaderPane
            chapter={currentChapter}
            {prevNav}
            {nextNav}
            onprev={previousChapter ? goToPreviousChapter : undefined}
            onnext={nextChapter ? goToNextChapter : undefined}
        >
            <ParagraphText verses={currentChapter.verses} />
        </ReaderPane>
    </div>
</div>
