<script lang="ts">
    import ChapterNavDivider from '@/components/reader/ChapterNavDivider.svelte';
    import ParagraphText from '@/components/reader/ParagraphText.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import ScrollSyncToggle from '@/components/reader/ScrollSyncToggle.svelte';
    import { bible, bookAbbrev, fetchChapter, loadBooks } from '@/lib/bible.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import {
        getNextChapter,
        getPreviousChapter,
        goToNextChapter,
        goToPreviousChapter,
        study,
    } from '@/lib/study.svelte.ts';
    import type { Chapter, ChapterNavTarget } from '@/lib/types/bible';
    import { verseNumbersInRange } from '@/lib/verseHighlight';
    import { createVerseHighlightScroller } from '@/lib/verseHighlightScroll';

    let leftScroll = $state<HTMLElement | null>(null);
    let rightScroll = $state<HTMLElement | null>(null);
    let syncing = $state(false);
    let suppressScrollSync = $state(false);
    let chapterA = $state<Chapter | null>(null);
    let chapterB = $state<Chapter | null>(null);
    let loading = $state(true);

    const highlightScroller = createVerseHighlightScroller();

    const previousChapter = $derived(getPreviousChapter());
    const nextChapter = $derived(getNextChapter());
    const readerStyle = $derived(getReaderStyle());

    const translationA = $derived(bible.translations.find((item) => item.id === study.translationId));
    const translationB = $derived(bible.translations.find((item) => item.id === study.translationBId));
    const highlightedVerses = $derived(
        study.verseHighlight
            ? verseNumbersInRange(study.verseHighlight.verse, study.verseHighlight.endVerse)
            : new Set<number>(),
    );

    $effect(() => {
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        const translationIdA = study.translationId;
        const translationIdB = study.translationBId;
        let cancelled = false;

        loading = true;
        chapterA = null;
        chapterB = null;

        Promise.all([
            loadBooks(translationIdA),
            fetchChapter(translationIdA, bookId, chapterNumber),
            fetchChapter(translationIdB, bookId, chapterNumber),
        ])
            .then(([, left, right]) => {
                if (! cancelled) {
                    chapterA = left;
                    chapterB = right;
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

    $effect(() => {
        if (! study.verseHighlight) {
            highlightScroller.reset();
        }
    });

    $effect(() => {
        const highlight = study.verseHighlight;

        if (! highlight || loading || ! chapterA) {
            return;
        }

        let cancelled = false;

        void (async () => {
            suppressScrollSync = study.scrollSync;

            await highlightScroller.scrollTo(
                [leftScroll, rightScroll],
                highlight,
                `${study.bookId}:${study.chapter}`,
            );

            if (! cancelled) {
                suppressScrollSync = false;
            }
        })();

        return () => {
            cancelled = true;
            suppressScrollSync = false;
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
        previousChapter
            ? chapterNavTarget(previousChapter.bookId, previousChapter.chapter)
            : undefined,
    );

    const nextNav = $derived(
        nextChapter ? chapterNavTarget(nextChapter.bookId, nextChapter.chapter) : undefined,
    );

    function syncScroll(source: HTMLElement, target: HTMLElement | null): void {
        if (! study.scrollSync || ! target || syncing || suppressScrollSync) {
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

    {#if loading || ! chapterA || ! chapterB}
        <div class="flex flex-1 items-center justify-center">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>
    {:else}
        <div class="grid min-h-0 flex-1 grid-cols-[1fr_auto_1fr]">
            <ReaderPane
                chapter={chapterA}
                translationAbbrev={translationA?.abbrev}
                bind:scrollRef={leftScroll}
                onscroll={handleLeftScroll}
            >
                <ParagraphText verses={chapterA.verses} {highlightedVerses} />
            </ReaderPane>

            <div class="border-base-300 flex w-12 flex-col items-center justify-between border-x py-4">
                {#if prevNav}
                    <ChapterNavDivider
                        direction="prev"
                        label={`${study.bookId != prevNav.bookId ? prevNav.bookAbbrev + ' ' : ''}${prevNav.chapter}`}
                        layout="compact"
                        onclick={goToPreviousChapter}
                    />
                {/if}
                {#if nextNav}
                    <ChapterNavDivider
                        direction="next"
                        label={`${study.bookId != nextNav.bookId ? nextNav.bookAbbrev + ' ' : ''}${nextNav.chapter}`}
                        layout="compact"
                        onclick={goToNextChapter}
                    />
                {/if}
            </div>

            <ReaderPane
                chapter={chapterB}
                translationAbbrev={translationB?.abbrev}
                bind:scrollRef={rightScroll}
                onscroll={handleRightScroll}
            >
                <ParagraphText verses={chapterB.verses} {highlightedVerses} />
            </ReaderPane>
        </div>
    {/if}
</div>
