<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { BibleService } from "../../mygo";
  import { setLocation, study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let refs = $state<
    { rank: number; bookId: string; bookName: string; chapter: number; verse: number; endVerse: number | null }[]
  >([]);
  let loading = $state(false);
  let error = $state("");

  $effect(() => {
    const bookId = study.bookId;
    const chapter = study.chapter;
    const verse = study.selectedVerse;
    loading = true;
    error = "";
    void BibleService.crossReferences(bookId, chapter, verse)
      .then((resp) => {
        refs = resp.references ?? [];
      })
      .catch((err: unknown) => {
        refs = [];
        error = err instanceof Error ? err.message : "Could not load cross references";
      })
      .finally(() => {
        loading = false;
      });
  });
</script>

<div class="flex h-full flex-col">
  <ColumnHeader contentType="cross-references" {slotIndex} showViewSelector />
  {#if loading && refs.length === 0}
    <p class="p-3 text-sm opacity-60">Loading…</p>
  {:else if error}
    <p class="p-3 text-sm text-error">{error}</p>
  {/if}
  <ul class="menu flex-1 overflow-y-auto p-2">
    {#each refs as ref (ref.rank)}
      <li>
        <button type="button" onclick={() => setLocation(ref.bookId, ref.chapter)}>
          <span class="badge badge-ghost badge-xs mr-2">{ref.rank}</span>
          {ref.bookName} {ref.chapter}:{ref.verse}{#if ref.endVerse && ref.endVerse !== ref.verse}–{ref.endVerse}{/if}
        </button>
      </li>
    {/each}
  </ul>
</div>
