<script lang="ts">
    import ScribePreviewModal from '@/components/scribe/ScribePreviewModal.svelte';
    import ScribeVerseField from '@/components/scribe/ScribeVerseField.svelte';
    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import {
        effectiveParagraphStart,
        entriesFromScribeVerses,
        fetchScribeDraft,
        scheduleScribeSave,
        serializeScribeDraft,
        type ScribeDraftEntry,
    } from '@/lib/scribe.svelte.ts';
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

    let entries = $state<Record<number, ScribeDraftEntry>>({});
    let loading = $state(true);
    let previewOpen = $state(false);

    const readerStyle = $derived(getReaderStyle());
    const verseNumbers = $derived(verses.map((verse) => verse.number));
    const hasParagraphOverrides = $derived(
        Object.values(entries).some((entry) => entry.paragraphStartOverride !== undefined),
    );

    $effect(() => {
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        let cancelled = false;

        loading = true;
        entries = {};

        void fetchScribeDraft(bookId, chapterNumber)
            .then((draft) => {
                if (! cancelled) {
                    entries = entriesFromScribeVerses(draft);
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

    function persist(): void {
        scheduleScribeSave(
            study.bookId,
            study.chapter,
            serializeScribeDraft(verseNumbers, entries),
        );
    }

    function updateVerse(verseNumber: number, value: string): void {
        entries = {
            ...entries,
            [verseNumber]: {
                ...entries[verseNumber],
                text: value,
            },
        };
        persist();
    }

    function toggleParagraphStart(verseNumber: number): void {
        if (verseNumber === 1) {
            return;
        }

        const source = verses.find((verse) => verse.number === verseNumber);
        const current = effectiveParagraphStart(
            verseNumber,
            source?.paragraphStart,
            entries[verseNumber]?.paragraphStartOverride,
        );
        const next = ! current;
        const sourceDefault = source?.paragraphStart ?? false;
        const entry = entries[verseNumber] ?? { text: '' };
        const updated: ScribeDraftEntry = { ...entry };

        if (next === sourceDefault) {
            delete updated.paragraphStartOverride;
        } else {
            updated.paragraphStartOverride = next;
        }

        entries = {
            ...entries,
            [verseNumber]: updated,
        };
        persist();
    }

    function resetParagraphBreaks(): void {
        entries = Object.fromEntries(
            Object.entries(entries).map(([verseNumber, entry]) => {
                const { paragraphStartOverride: _, ...rest } = entry;

                return [Number(verseNumber), rest];
            }),
        );
        persist();
    }
</script>

<div class="flex h-full min-h-0 flex-col overflow-y-auto px-4 py-4" style={readerStyle}>
    <div class="flex justify-between">
        <ChapterHeading title="{book} {chapter}" />

        <div class="flex flex-wrap gap-2">
            <button type="button" class="btn btn-ghost btn-sm" disabled={loading} onclick={() => (previewOpen = true)}>
                Preview
            </button>
            {#if hasParagraphOverrides}
                <button type="button" class="btn btn-ghost btn-sm" onclick={resetParagraphBreaks}>
                    Reset paragraph breaks
                </button>
            {/if}
        </div>
    </div>


    {#if loading}
        <div class="flex flex-1 items-center justify-center">
            <span class="loading loading-spinner loading-md text-primary"></span>
        </div>
    {:else}
        {#each verses as verse (verse.number)}
            <ScribeVerseField
                {verse}
                value={entries[verse.number]?.text ?? ''}
                paragraphStart={effectiveParagraphStart(
                    verse.number,
                    verse.paragraphStart,
                    entries[verse.number]?.paragraphStartOverride,
                )}
                onupdate={(value) => updateVerse(verse.number, value)}
                onToggleParagraph={() => toggleParagraphStart(verse.number)}
            />
        {/each}
    {/if}
</div>

<ScribePreviewModal
    open={previewOpen}
    title="{book} {chapter}"
    sourceVerses={verses}
    {entries}
    onclose={() => (previewOpen = false)}
/>
