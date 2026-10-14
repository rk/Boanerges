<script lang="ts">
  import { groupVersesIntoParagraphs } from "../../lib/paragraphs";
  import type { Verse } from "../../mygo";
  import VerseSpan from "./VerseSpan.svelte";

  let {
    verses,
    bookId,
    chapter,
    showContextMenu = false,
  }: {
    verses: Verse[];
    bookId: string;
    chapter: number;
    showContextMenu?: boolean;
  } = $props();

  const paragraphs = $derived(groupVersesIntoParagraphs(verses));
</script>

<div class="reader-prose space-y-4">
  {#each paragraphs as paragraph (paragraph[0]?.number ?? 0)}
    <p>
      {#each paragraph as verse, index (verse.number)}
        {#if index > 0}
          {" "}
        {/if}
        <VerseSpan {verse} {bookId} {chapter} {showContextMenu} />
      {/each}
    </p>
  {/each}
</div>
