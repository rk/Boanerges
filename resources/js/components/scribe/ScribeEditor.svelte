<script lang="ts">
    import ScribeVerseField from '@/components/scribe/ScribeVerseField.svelte';
    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import { loadScribeDraft, scheduleScribeSave } from '@/lib/scribe.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { study } from '@/lib/study.svelte.ts';
    import type { Verse } from '@/lib/types/bible';

    let {
        book,
        chapter,
        verses,
    }: {
        book: string;
        chapter: number;
        verses: Verse[];
    } = $props();

    let draft = $state<Record<number, string>>({});

    const readerStyle = $derived(getReaderStyle());

    $effect(() => {
        draft = loadScribeDraft(study.bookId, study.chapter);
    });

    function updateVerse(verseNumber: number, value: string): void {
        draft = { ...draft, [verseNumber]: value };
        scheduleScribeSave(study.bookId, study.chapter, draft);
    }
</script>

<div class="flex h-full min-h-0 flex-col overflow-y-auto px-4 py-4" style={readerStyle}>
    <ChapterHeading title="{book} {chapter}" />

    {#each verses as verse (verse.number)}
        <ScribeVerseField
            {verse}
            value={draft[verse.number] ?? ''}
            onupdate={(value) => updateVerse(verse.number, value)}
        />
    {/each}
</div>
