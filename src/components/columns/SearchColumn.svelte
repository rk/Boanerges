<script lang="ts">
  import { BibleService } from "../../mygo";
  import { study } from "../../lib/study.svelte.ts";

  let query = $state("");
  let results = $state<{ bookId: string; chapter: number; verse: number; snippet: string }[]>([]);
  let offset = $state(0);
  let hasMore = $state(false);
  const pageSize = 20;

  async function runSearch(reset = true) {
    if (query.trim().length < 2) {
      results = [];
      hasMore = false;
      return;
    }
    if (reset) offset = 0;
    const resp = await BibleService.search(query, study.translationId, pageSize, offset);
    results = reset ? resp.results : [...results, ...resp.results];
    hasMore = resp.hasMore;
  }

  async function loadMore() {
    offset += pageSize;
    await runSearch(false);
  }
</script>

<div class="flex h-full flex-col" data-testid="search-column">
  <div class="border-b border-base-300 p-3">
    <input
      class="input input-bordered input-sm w-full"
      placeholder="Search…"
      data-testid="search-input"
      bind:value={query}
      onkeydown={(e) => e.key === "Enter" && runSearch()}
    />
  </div>
  <ul class="menu flex-1 overflow-y-auto p-2">
    {#each results as hit}
      <li>
        <button
          type="button"
          class="text-left whitespace-normal"
          data-testid="search-hit"
          onclick={() => {
            study.bookId = hit.bookId;
            study.chapter = hit.chapter;
            study.selectedVerse = hit.verse;
          }}
        >
          {@html hit.snippet}
        </button>
      </li>
    {/each}
  </ul>
  {#if hasMore}
    <div class="border-t border-base-300 p-2">
      <button type="button" class="btn btn-sm btn-block" data-testid="search-load-more" onclick={loadMore}>Load more</button>
    </div>
  {/if}
</div>
