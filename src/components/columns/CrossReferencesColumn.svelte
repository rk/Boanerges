<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { BibleService } from "../../mygo";
  import { crossReferenceKey, normalizeCrossReferences, type CrossReferenceHit } from "../../lib/crossReferences";
  import { canonBookName } from "../../lib/canonBookNames";
  import { navigateToVerse, study } from "../../lib/study.svelte.ts";

  const origin = $derived(study.crossRefsTarget);

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let refs = $state<CrossReferenceHit[]>([]);
  let loading = $state(false);
  let error = $state("");

  $effect(() => {
    study.crossRefsTick;
    const target = study.crossRefsTarget;
    if (!target) {
      refs = [];
      loading = false;
      error = "";
      return;
    }
    const { bookId, chapter, verse } = target;
    loading = true;
    error = "";
    refs = [];
    let cancelled = false;
    void BibleService.crossReferences(bookId, chapter, verse)
      .then((resp) => {
        if (cancelled) {
          return;
        }
        refs = normalizeCrossReferences(resp);
      })
      .catch((err: unknown) => {
        if (cancelled) {
          return;
        }
        refs = [];
        error = err instanceof Error ? err.message : "Could not load cross references";
      })
      .finally(() => {
        if (!cancelled) {
          loading = false;
        }
      });
    return () => {
      cancelled = true;
    };
  });
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" data-testid="cross-references-column">
  <ColumnHeader contentType="cross-references" {slotIndex} showViewSelector />
  {#if !study.crossRefsTarget}
    <p class="p-3 text-sm opacity-60">Right-click a verse in the primary Bible for cross-references.</p>
  {:else if loading}
    <p class="p-3 text-sm opacity-60">Loading…</p>
  {:else if error}
    <p class="p-3 text-sm text-error">{error}</p>
  {:else if refs.length === 0}
    <p class="p-3 text-sm opacity-60">No cross-references for this verse.</p>
  {/if}
  <div class="flex min-h-0 flex-1 flex-col overflow-y-auto p-2" role="list">
    {#if origin && !loading && !error}
      <button
        type="button"
        class="hover:bg-base-200 w-full rounded-btn px-3 py-2 text-left text-sm font-medium whitespace-normal"
        data-testid="cross-ref-origin"
        role="listitem"
        onclick={() => navigateToVerse(origin.bookId, origin.chapter, origin.verse)}
      >
        <span class="badge badge-primary badge-xs mr-2 shrink-0">↩</span>
        {canonBookName(origin.bookId)} {origin.chapter}:{origin.verse}
      </button>
    {/if}
    {#each refs as ref, index (crossReferenceKey(ref, index))}
      <button
        type="button"
        class="hover:bg-base-200 w-full rounded-btn px-3 py-2 text-left text-sm whitespace-normal"
        data-testid="cross-ref-hit"
        role="listitem"
        onclick={() => navigateToVerse(ref.bookId, ref.chapter, ref.verse)}
      >
        <span class="badge badge-neutral badge-xs mr-2 shrink-0">{ref.rank}</span>
        {ref.bookName} {ref.chapter}:{ref.verse}{#if ref.endVerse !== null && ref.endVerse !== ref.verse}–{ref.endVerse}{/if}
      </button>
    {/each}
  </div>
</div>
