<script lang="ts">
  import { tick } from "svelte";
  import ScribePreviewModal from "./ScribePreviewModal.svelte";
  import ScribeVerseSpan from "./ScribeVerseSpan.svelte";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import {
    applyPrimaryParagraphBreaks,
    effectiveParagraphStart,
    entriesFromScribeVerses,
    fetchScribeDraft,
    resetInterVerseParagraphBreaks,
    scheduleScribeSave,
    serializeScribeDraft,
    type ScribeDraftEntry,
  } from "../../lib/scribe.svelte.ts";
  import { study } from "../../lib/study.svelte.ts";
  import type { Verse } from "../../mygo";

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
  let spanComponents = $state<
    Record<number, { setText: (text: string, options?: { force?: boolean }) => void } | undefined>
  >({});

  const verseNumbers = $derived(verses.map((verse) => verse.number));
  const documentKey = $derived(`${study.bookId}-${study.chapter}`);

  $effect(() => {
    const bookId = study.bookId;
    const chapterNumber = study.chapter;
    let cancelled = false;

    loading = true;
    entries = {};
    spanComponents = {};

    void fetchScribeDraft(bookId, chapterNumber)
      .then(async (draft) => {
        if (!cancelled) {
          entries = entriesFromScribeVerses(draft);
          loading = false;
          await tick();
          hydrateSpans(true);
        }
      })
      .catch(async () => {
        if (!cancelled) {
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
      spanComponents[verseNumber]?.setText(entries[verseNumber]?.text ?? "", { force });
    }
  }

  function persist(): void {
    scheduleScribeSave(study.bookId, study.chapter, serializeScribeDraft(verseNumbers, entries));
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

    const current = effectiveParagraphStart(verseNumber, entries[verseNumber]?.paragraphStartOverride);
    const next = !current;
    const entry = entries[verseNumber] ?? { text: "" };
    const updated: ScribeDraftEntry = { ...entry };

    if (next === false) {
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

  export function openPreview(): void {
    previewOpen = true;
  }

  export function setLineBreaksFromPrimary(): void {
    entries = applyPrimaryParagraphBreaks(verseNumbers, entries, verses);
    persist();
  }

  export function resetLineBreaks(): void {
    entries = resetInterVerseParagraphBreaks(verseNumbers, entries);
    persist();
  }
</script>

<div
  class="flex h-full min-h-0 flex-col overflow-y-auto border-x border-base-300 px-4 py-4"
  style={readerStyle()}
>
  {#if loading}
    <div class="flex flex-1 items-center justify-center">
      <span class="loading loading-md loading-spinner text-primary"></span>
    </div>
  {:else}
    {#key documentKey}
      <div
        class="scribe-document prose max-w-none"
        role="region"
        aria-label="Scribe draft"
        data-testid="scribe-editor"
      >
        {#each verses as verse (verse.number)}
          <ScribeVerseSpan
            bind:this={spanComponents[verse.number]}
            verseNumber={verse.number}
            paragraphStart={effectiveParagraphStart(
              verse.number,
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
