<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { BibleService } from "../../mygo";
  import { formatScriptureReference } from "../../lib/scriptureReference";
  import { navigateToVerse, study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let query = $state("");
  let results = $state<{ bookId: string; chapter: number; verse: number; snippet: string }[]>([]);
  let offset = $state(0);
  let hasMore = $state(false);
  let searching = $state(false);
  const pageSize = 20;
  let debounce: ReturnType<typeof setTimeout> | null = null;

  async function runSearch(reset = true) {
    if (query.trim().length < 2) {
      results = [];
      hasMore = false;
      searching = false;
      return;
    }
    searching = true;
    if (reset) offset = 0;
    const resp = await BibleService.search(query, study.translationId, pageSize, offset);
    results = reset ? resp.results : [...results, ...resp.results];
    hasMore = resp.hasMore;
    searching = false;
  }

  $effect(() => {
    const q = query;
    if (debounce) clearTimeout(debounce);
    debounce = setTimeout(() => {
      void runSearch(true);
    }, 300);
    return () => {
      if (debounce) clearTimeout(debounce);
    };
  });

  async function loadMore() {
    offset += pageSize;
    await runSearch(false);
  }

  function goToHit(hit: { bookId: string; chapter: number; verse: number }) {
    navigateToVerse(hit.bookId, hit.chapter, hit.verse);
  }
</script>

<div class="flex h-full min-h-0 flex-col" data-testid="search-column">
  <ColumnHeader contentType="search" {slotIndex} showViewSelector>
    <input
      class="input input-bordered input-sm w-full"
      placeholder="Search…"
      data-testid="search-input"
      bind:value={query}
      onkeydown={(e) => e.key === "Enter" && runSearch()}
    />
  </ColumnHeader>
  {#if searching && results.length === 0}
    <p class="p-3 text-sm opacity-60">Searching…</p>
  {/if}
  <ul class="menu min-h-0 flex-1 overflow-y-auto p-2">
    {#each results as hit}
      <li>
        <button type="button" class="text-left whitespace-normal" data-testid="search-hit" onclick={() => goToHit(hit)}>
          <span class="text-primary block text-xs font-medium">
            {formatScriptureReference(hit.bookId, hit.chapter, hit.verse)}
          </span>
          <span>{@html hit.snippet}</span>
        </button>
      </li>
    {/each}
  </ul>
  {#if hasMore}
    <div class="border-t border-base-300 shrink-0 p-2">
      <button type="button" class="btn btn-sm btn-block" data-testid="search-load-more" onclick={loadMore}>Load more</button>
    </div>
  {/if}
</div>
