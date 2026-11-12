<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    onclick,
    testId,
    body,
    lead,
  }: {
    onclick: () => void;
    testId?: string;
    body: Snippet;
    lead?: Snippet;
  } = $props();
</script>

<button
  type="button"
  class="hover:bg-base-200 w-full rounded-btn px-3 py-2 text-left text-sm whitespace-normal"
  data-testid={testId}
  role="listitem"
  {onclick}
>
  {#if lead}
    <span class="mb-1 flex items-center gap-2">{@render lead()}</span>
  {/if}
  {@render body()}
</button>
