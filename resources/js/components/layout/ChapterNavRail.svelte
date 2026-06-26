<script lang="ts">
    import ChapterNavDivider from '@/components/reader/ChapterNavDivider.svelte';
    import { bookAbbrev } from '@/lib/bible.svelte.ts';
    import {
        getNextChapter,
        getPreviousChapter,
        goToNextChapter,
        goToPreviousChapter,
        study,
    } from '@/lib/study.svelte.ts';

    const previousChapter = $derived(getPreviousChapter());
    const nextChapter = $derived(getNextChapter());

    const prevNav = $derived(
        previousChapter
            ? {
                bookId: previousChapter.bookId,
                bookAbbrev: bookAbbrev(previousChapter.bookId),
                chapter: previousChapter.chapter,
            }
            : undefined,
    );

    const nextNav = $derived(
        nextChapter
            ? {
                bookId: nextChapter.bookId,
                bookAbbrev: bookAbbrev(nextChapter.bookId),
                chapter: nextChapter.chapter,
            }
            : undefined,
    );
</script>

<aside
    class="border-base-300 flex w-12 shrink-0 flex-col items-center justify-between border-r py-4"
    aria-label="Chapter navigation"
>
    {#if prevNav}
        <ChapterNavDivider
            direction="prev"
            label={`${study.bookId !== prevNav.bookId ? `${prevNav.bookAbbrev} ` : ''}${prevNav.chapter}`}
            layout="compact"
            onclick={goToPreviousChapter}
        />
    {/if}
    {#if nextNav}
        <ChapterNavDivider
            direction="next"
            label={`${study.bookId !== nextNav.bookId ? `${nextNav.bookAbbrev} ` : ''}${nextNav.chapter}`}
            layout="compact"
            onclick={goToNextChapter}
        />
    {/if}
</aside>
