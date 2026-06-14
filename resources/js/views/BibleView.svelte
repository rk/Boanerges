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
        navLabel,
    } from '@/lib/study.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';

    const currentChapter = $derived(getCurrentChapter());
    const previousChapter = $derived(getPreviousChapter());
    const nextChapter = $derived(getNextChapter());
    const readerStyle = $derived(getReaderStyle());

    const prevLabel = $derived(
        previousChapter
            ? navLabel(getChapter(previousChapter.bookId, previousChapter.chapter).bookAbbrev, previousChapter.chapter, 'prev')
            : undefined,
    );

    const nextLabel = $derived(
        nextChapter
            ? navLabel(getChapter(nextChapter.bookId, nextChapter.chapter).bookAbbrev, nextChapter.chapter, 'next')
            : undefined,
    );
</script>

<div class="h-full min-h-0" style={readerStyle}>
    <div class="mx-auto h-full max-w-prose">
        <ReaderPane
            chapter={currentChapter}
            prevLabel={prevLabel}
            nextLabel={nextLabel}
            onprev={previousChapter ? goToPreviousChapter : undefined}
            onnext={nextChapter ? goToNextChapter : undefined}
        >
            <ParagraphText verses={currentChapter.verses} />
        </ReaderPane>
    </div>
</div>
