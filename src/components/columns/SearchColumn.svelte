<script lang="ts">
  import { BibleService } from "../../mygo";
  import { study } from "../../lib/study.svelte.ts";

  let query = $state("");
  let results = $state<{ bookId: string; chapter: number; verse: number; snippet: string }[]>([]);

  async function runSearch() {
    if (query.trim().length < 2) {
      results = [];
      return;
    }
    const resp = await BibleService.search(query, study.translationId, 50, 0);
    results = resp.results;
  }
</script>

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 p-3">
    <input class="input input-bordered input-sm w-full" placeholder="Search…" bind:value={query} onkeydown={(e) => e.key === "Enter" && runSearch()} />
  </div>
  <ul class="menu flex-1 overflow-y-auto p-2">
    {#each results as hit}
      <li><button type="button" class="text-left whitespace-normal" onclick={() => { study.bookId = hit.bookId; study.chapter = hit.chapter; study.selectedVerse = hit.verse; }}>{@html hit.snippet}</button></li>
    {/each}
  </ul>
</div>
