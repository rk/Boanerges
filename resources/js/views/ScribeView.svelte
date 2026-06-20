<script lang="ts">
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import VerseText from '@/components/reader/VerseText.svelte';
    import ScribeEditor from '@/components/scribe/ScribeEditor.svelte';
    import { bible, fetchChapter, loadBooks } from '@/lib/bible.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { study } from '@/lib/study.svelte.ts';
    import type { Chapter } from '@/lib/types/bible';
    import { verseNumbersInRange } from '@/lib/verseHighlight';
    import { createVerseHighlightScroller } from '@/lib/verseHighlightScroll';

    let chapterA = $state<Chapter | null>(null);
    let chapterB = $state<Chapter | null>(null);
    let loading = $state(true);
    let leftScroll = $state<HTMLElement | null>(null);
    let rightScroll = $state<HTMLElement | null>(null);

    const highlightScroller = createVerseHighlightScroller();
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

        void highlightScroller.scrollTo(
            [leftScroll, rightScroll],
            highlight,
            `${study.bookId}:${study.chapter}`,
        );
    });
</script>

<div class="grid h-full min-h-0 grid-cols-3" style={readerStyle}>
    {#if loading || ! chapterA || ! chapterB}
        <div class="col-span-3 flex items-center justify-center">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>
    {:else}
        <ReaderPane chapter={chapterA} translationAbbrev={translationA?.abbrev} bind:scrollRef={leftScroll}>
            <VerseText verses={chapterA.verses} {highlightedVerses} />
        </ReaderPane>

        <ScribeEditor book={chapterA.book} chapter={chapterA.chapter} verses={chapterA.verses} />

        <ReaderPane chapter={chapterB} translationAbbrev={translationB?.abbrev} bind:scrollRef={rightScroll}>
            <VerseText verses={chapterB.verses} {highlightedVerses} />
        </ReaderPane>
    {/if}
</div>
