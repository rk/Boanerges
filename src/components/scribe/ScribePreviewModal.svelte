<script lang="ts">
  import ScribeParagraphText from "./ScribeParagraphText.svelte";
  import { toPreviewVerses } from "../../lib/paragraphs";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { effectiveParagraphStart, type ScribeDraftEntry } from "../../lib/scribe.svelte.ts";
  import type { Verse } from "../../mygo";

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

  let dialog = $state<HTMLDialogElement | undefined>();

  const previewVerses = $derived(
    toPreviewVerses(
      sourceVerses,
      (verseNumber) => entries[verseNumber]?.text ?? "",
      (verseNumber) => effectiveParagraphStart(verseNumber, entries[verseNumber]?.paragraphStartOverride),
    ),
  );

  const hasContent = $derived(previewVerses.length > 0);

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
    <h2 class="text-lg font-semibold">{title}</h2>

    <div class="min-h-0 flex-1 overflow-y-auto py-2" style={readerStyle()}>
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
