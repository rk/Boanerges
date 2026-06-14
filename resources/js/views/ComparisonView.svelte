<script lang="ts">
    import ParagraphText from '@/components/reader/ParagraphText.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import ScrollSyncToggle from '@/components/reader/ScrollSyncToggle.svelte';
    import ChapterNavDivider from '@/components/reader/ChapterNavDivider.svelte';
    import { getChapter, translations } from '@/lib/mock/chapter';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import {
        getCurrentChapter,
        getNextChapter,
        getPreviousChapter,
        goToNextChapter,
        goToPreviousChapter,
        navLabel,
        study,
    } from '@/lib/study.svelte.ts';

    let leftScroll: HTMLElement | null = $state(null);
    let rightScroll: HTMLElement | null = $state(null);
    let syncing = $state(false);

    const currentChapter = $derived(getCurrentChapter());
    const previousChapter = $derived(getPreviousChapter());
    const nextChapter = $derived(getNextChapter());
    const readerStyle = $derived(getReaderStyle());

    const translationA = $derived(translations.find((item) => item.id === study.translationId));
    const translationB = $derived(translations.find((item) => item.id === study.translationBId));

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

    function syncScroll(source: HTMLElement, target: HTMLElement | null): void {
        if (! study.scrollSync || ! target || syncing) {
            return;
        }

        syncing = true;

        const sourceRange = source.scrollHeight - source.clientHeight;
        const targetRange = target.scrollHeight - target.clientHeight;
        const ratio = sourceRange > 0 ? source.scrollTop / sourceRange : 0;

        target.scrollTop = ratio * targetRange;

        syncing = false;
    }

    function handleLeftScroll(): void {
        if (leftScroll) {
            syncScroll(leftScroll, rightScroll);
        }
    }

    function handleRightScroll(): void {
        if (rightScroll) {
            syncScroll(rightScroll, leftScroll);
        }
    }
</script>

<div class="flex h-full min-h-0 flex-col" style={readerStyle}>
    <div class="border-base-300 flex shrink-0 items-center justify-end border-b px-4 py-2">
        <ScrollSyncToggle />
    </div>

    <div class="grid min-h-0 flex-1 grid-cols-[1fr_auto_1fr]">
        <ReaderPane
            chapter={currentChapter}
            translationAbbrev={translationA?.abbrev}
            bind:scrollRef={leftScroll}
            onscroll={handleLeftScroll}
        >
            <ParagraphText verses={currentChapter.verses} />
        </ReaderPane>

        <div class="border-base-300 flex w-12 flex-col items-center justify-between border-x py-4">
            {#if prevLabel}
                <ChapterNavDivider
                    label={prevLabel}
                    direction="vertical"
                    onclick={goToPreviousChapter}
                />
            {/if}
            <div class="divider divider-vertical flex-1"></div>
            {#if nextLabel}
                <ChapterNavDivider
                    label={nextLabel}
                    direction="vertical"
                    onclick={goToNextChapter}
                />
            {/if}
        </div>

        <ReaderPane
            chapter={currentChapter}
            translationAbbrev={translationB?.abbrev}
            bind:scrollRef={rightScroll}
            onscroll={handleRightScroll}
        >
            <ParagraphText verses={currentChapter.verses} />
        </ReaderPane>
    </div>
</div>
