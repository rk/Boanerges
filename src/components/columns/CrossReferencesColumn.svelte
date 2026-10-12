<script lang="ts">
  import { BibleService } from "../../mygo";
  import { setLocation, study } from "../../lib/study.svelte.ts";

  let refs = $state<{ rank: number; bookId: string; bookName: string; chapter: number; verse: number; endVerse: number | null }[]>([]);

  $effect(() => {
    void BibleService.crossReferences(study.bookId, study.chapter, study.selectedVerse).then((resp) => {
      refs = resp.references ?? [];
    });
  });
</script>

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">Cross references · v{study.selectedVerse}</div>
  <ul class="menu flex-1 overflow-y-auto p-2">
    {#each refs as ref}
      <li>
        <button type="button" onclick={() => setLocation(ref.bookId, ref.chapter)}>
          {ref.bookName} {ref.chapter}:{ref.verse}{#if ref.endVerse && ref.endVerse !== ref.verse}–{ref.endVerse}{/if}
        </button>
      </li>
    {/each}
  </ul>
</div>
