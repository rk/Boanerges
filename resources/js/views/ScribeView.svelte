<script lang="ts">
    import ScribeEditor from '@/components/scribe/ScribeEditor.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import VerseText from '@/components/reader/VerseText.svelte';
    import { bible, fetchChapter, loadBooks } from '@/lib/bible.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { study } from '@/lib/study.svelte.ts';
    import type { Chapter } from '@/lib/types/bible';

    let chapterA = $state<Chapter | null>(null);
    let chapterB = $state<Chapter | null>(null);
    let loading = $state(true);

    const readerStyle = $derived(getReaderStyle());

    const translationA = $derived(bible.translations.find((item) => item.id === study.translationId));
    const translationB = $derived(bible.translations.find((item) => item.id === study.translationBId));

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
</script>

<div class="grid h-full min-h-0 grid-cols-3" style={readerStyle}>
    {#if loading || ! chapterA || ! chapterB}
        <div class="col-span-3 flex items-center justify-center">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>
    {:else}
        <ReaderPane chapter={chapterA} translationAbbrev={translationA?.abbrev}>
            <VerseText verses={chapterA.verses} />
        </ReaderPane>

        <ScribeEditor book={chapterA.book} chapter={chapterA.chapter} verses={chapterA.verses} />

        <ReaderPane chapter={chapterB} translationAbbrev={translationB?.abbrev}>
            <VerseText verses={chapterB.verses} />
        </ReaderPane>
    {/if}
</div>
