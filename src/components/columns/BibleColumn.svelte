<script lang="ts">
  import { BibleService, type Chapter } from "../../mygo";
  import VerseLine from "../reader/VerseLine.svelte";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { study } from "../../lib/study.svelte.ts";

  let {
    translationId,
    primary = false,
    slotIndex = 0,
    scrollRef = $bindable<HTMLElement | null>(null),
    onscroll,
  }: {
    translationId: string;
    primary?: boolean;
    slotIndex?: number;
    columnType?: string;
    scrollRef?: HTMLElement | null;
    onscroll?: () => void;
  } = $props();

  let chapter = $state<Chapter | null>(null);
  let loading = $state(true);

  $effect(() => {
    loading = true;
    void BibleService.chapter(translationId, study.bookId, study.chapter).then((resp) => {
      chapter = resp.chapter;
      loading = false;
    });
  });
</script>

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">
    {primary ? "Primary" : `Translation ${slotIndex + 1}`} · {translationId.toUpperCase()}
  </div>
  <article
    bind:this={scrollRef}
    class="prose max-w-none flex-1 overflow-y-auto p-4"
    style={readerStyle()}
    data-testid={primary ? "bible-reader-primary" : "bible-reader-secondary"}
    onscroll={onscroll}
  >
    {#if loading}
      <p>Loading…</p>
    {:else if chapter}
      <h2 data-testid="chapter-heading">{chapter.book} {chapter.chapter}</h2>
      {#each chapter.verses as verse (verse.number)}
        <VerseLine {verse} highlighted={primary && verse.number === study.selectedVerse} />
      {/each}
    {/if}
  </article>
</div>
