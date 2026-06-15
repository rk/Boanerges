<script lang="ts">
    import { tick } from 'svelte';

    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import ScribePreviewModal from '@/components/scribe/ScribePreviewModal.svelte';
    import ScribeVerseSpan from '@/components/scribe/ScribeVerseSpan.svelte';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import {
        effectiveParagraphStart,
        entriesFromScribeVerses,
        fetchScribeDraft,
        scheduleScribeSave,
        serializeScribeDraft

    } from '@/lib/scribe.svelte.ts';
import type {ScribeDraftEntry} from '@/lib/scribe.svelte.ts';
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
    let spanComponents = $state<Record<number, { setText: (text: string, options?: { force?: boolean }) => void } | undefined>>({});

    const readerStyle = $derived(getReaderStyle());
    const verseNumbers = $derived(verses.map((verse) => verse.number));
    const documentKey = $derived(`${study.bookId}-${study.chapter}`);
    const hasParagraphOverrides = $derived(
        Object.values(entries).some((entry) => entry.paragraphStartOverride !== undefined),
    );

    $effect(() => {
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        let cancelled = false;

        loading = true;
        entries = {};
        spanComponents = {};

        void fetchScribeDraft(bookId, chapterNumber)
            .then(async (draft) => {
                if (! cancelled) {
                    entries = entriesFromScribeVerses(draft);
                    loading = false;
                    await tick();
                    hydrateSpans(true);
                }
            })
            .catch(async () => {
                if (! cancelled) {
                    loading = false;
                    await tick();
                    hydrateSpans(true);
                }
            });

        return () => {
            cancelled = true;
        };
    });

    function hydrateSpans(force = false): void {
        for (const verseNumber of verseNumbers) {
            spanComponents[verseNumber]?.setText(entries[verseNumber]?.text ?? '', { force });
        }
    }

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

<div class="flex h-full min-h-0 flex-col overflow-y-auto px-4 py-4 border-base-300 border-x" style={readerStyle}>
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
        {#key documentKey}
            <div class="scribe-document reader-prose" role="region" aria-label="Scribe draft">
                {#each verses as verse (verse.number)}
                    <ScribeVerseSpan
                        bind:this={spanComponents[verse.number]}
                        verseNumber={verse.number}
                        paragraphStart={effectiveParagraphStart(
                            verse.number,
                            verse.paragraphStart,
                            entries[verse.number]?.paragraphStartOverride,
                        )}
                        oninput={(value) => updateVerse(verse.number, value)}
                        onToggleParagraph={() => toggleParagraphStart(verse.number)}
                    />
                {/each}
            </div>
        {/key}
    {/if}
</div>

<ScribePreviewModal
    open={previewOpen}
    title="{book} {chapter}"
    sourceVerses={verses}
    {entries}
    onclose={() => (previewOpen = false)}
/>
