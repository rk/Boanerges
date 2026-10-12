<script lang="ts">
  import { BibleService, type Chapter } from "../../mygo";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { study } from "../../lib/study.svelte.ts";

  let {
    scrollRefs = $bindable<[HTMLElement | null, HTMLElement | null]>([null, null]),
    onscroll,
  }: {
    scrollRefs?: [HTMLElement | null, HTMLElement | null];
    onscroll?: () => void;
  } = $props();

  let left = $state<Chapter | null>(null);
  let right = $state<Chapter | null>(null);
  let leftEl = $state<HTMLElement | null>(null);
  let rightEl = $state<HTMLElement | null>(null);

  $effect(() => {
    scrollRefs = [leftEl, rightEl];
  });

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

<div class="flex h-full flex-col" data-testid="comparison-column">
  <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">Comparison</div>
  <div class="grid flex-1 min-h-0 grid-cols-2 divide-x divide-base-300 overflow-hidden">
    <article bind:this={leftEl} class="overflow-y-auto p-3 text-sm" style={readerStyle()} data-testid="comparison-pane" onscroll={onscroll}>
      <h3 class="mb-2 font-semibold">{study.translationBId.toUpperCase()}</h3>
      {#if left}
        {#each left.verses as verse}
          <p class="mb-2"><sup>{verse.number}</sup> {@html verse.text}</p>
        {/each}
      {/if}
    </article>
    <article bind:this={rightEl} class="overflow-y-auto p-3 text-sm" style={readerStyle()} data-testid="comparison-pane" onscroll={onscroll}>
      <h3 class="mb-2 font-semibold">{study.translationCId.toUpperCase()}</h3>
      {#if right}
        {#each right.verses as verse}
          <p class="mb-2"><sup>{verse.number}</sup> {@html verse.text}</p>
        {/each}
      {/if}
    </article>
  </div>
</div>
