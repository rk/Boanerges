<script lang="ts">
    import ScribeVerseField from '@/components/scribe/ScribeVerseField.svelte';
    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import { loadScribeDraft, scheduleScribeSave } from '@/lib/scribe.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { getCurrentChapter, study } from '@/lib/study.svelte.ts';

    let draft = $state<Record<number, string>>({});

    const currentChapter = $derived(getCurrentChapter());
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
    <ChapterHeading title="{currentChapter.book} {currentChapter.chapter}" />

    {#each currentChapter.verses as verse (verse.number)}
        <ScribeVerseField
            {verse}
            value={draft[verse.number] ?? ''}
            onupdate={(value) => updateVerse(verse.number, value)}
        />
    {/each}
</div>
