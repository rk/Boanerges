<script lang="ts">
  import { BibleService, type Chapter } from "../../mygo";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { study } from "../../lib/study.svelte.ts";

  let left = $state<Chapter | null>(null);
  let right = $state<Chapter | null>(null);

  $effect(() => {
    void Promise.all([
      BibleService.chapter(study.translationBId, study.bookId, study.chapter),
      BibleService.chapter(study.translationCId, study.bookId, study.chapter),
    ]).then(([a, b]) => {
      left = a.chapter;
      right = b.chapter;
    });
  });
</script>

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">Comparison</div>
  <div class="grid flex-1 min-h-0 grid-cols-2 divide-x divide-base-300 overflow-hidden">
    {#each [{ label: study.translationBId, data: left }, { label: study.translationCId, data: right }] as col}
      <article class="overflow-y-auto p-3 text-sm" style={readerStyle()}>
        <h3 class="mb-2 font-semibold">{col.label.toUpperCase()}</h3>
        {#if col.data}
          {#each col.data.verses as verse}
            <p class="mb-2"><sup>{verse.number}</sup> {@html verse.text}</p>
          {/each}
        {/if}
      </article>
    {/each}
  </div>
</div>
