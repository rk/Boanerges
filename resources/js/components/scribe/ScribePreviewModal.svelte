<script lang="ts">
    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import ScribeParagraphText from '@/components/scribe/ScribeParagraphText.svelte';
    import { toPreviewVerses } from '@/lib/paragraphs.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { effectiveParagraphStart } from '@/lib/scribe.svelte.ts';
    import type { ScribeDraftEntry } from '@/lib/scribe.svelte.ts';
    import type { Verse } from '@/lib/types/bible';

    let {
        open = false,
        title,
        sourceVerses,
        entries,
        onclose,
    }: {
        open?: boolean;
        title: string;
        sourceVerses: Verse[];
        entries: Record<number, ScribeDraftEntry>;
        onclose?: () => void;
    } = $props();

    let dialog: HTMLDialogElement | undefined = $state();

    const previewVerses = $derived(
        toPreviewVerses(
            sourceVerses,
            (verseNumber) => entries[verseNumber]?.text ?? '',
            (verseNumber) => {
                const source = sourceVerses.find((verse) => verse.number === verseNumber);

                return effectiveParagraphStart(
                    verseNumber,
                    source?.paragraphStart,
                    entries[verseNumber]?.paragraphStartOverride,
                );
            },
        ),
    );

    const hasContent = $derived(previewVerses.length > 0);

    const readerStyle = $derived(getReaderStyle());

    $effect(() => {
        if (open) {
            dialog?.showModal();
        } else {
            dialog?.close();
        }
    });

    function handleClose(): void {
        onclose?.();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box flex max-h-[85vh] max-w-3xl flex-col">
        <ChapterHeading {title} />

        <div class="min-h-0 flex-1 overflow-y-auto py-2" style={readerStyle}>
            {#if hasContent}
                <ScribeParagraphText verses={previewVerses} />
            {:else}
                <p class="text-base-content/70 py-8 text-center text-sm">Nothing to preview yet.</p>
            {/if}
        </div>

        <div class="modal-action">
            <form method="dialog">
                <button type="submit" class="btn">Close</button>
            </form>
        </div>
    </div>

    <form method="dialog" class="modal-backdrop">
        <button type="submit" aria-label="Close preview">Close</button>
    </form>
</dialog>
