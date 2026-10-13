<script lang="ts">
  import FormattedVerseText from "../reader/FormattedVerseText.svelte";
  import { buildScribePreviewParagraphs } from "../../lib/paragraphs";
  import type { Verse } from "../../mygo";

  let { verses }: { verses: Verse[] } = $props();

  const paragraphs = $derived(buildScribePreviewParagraphs(verses));
</script>

<div class="prose max-w-none space-y-4">
  {#each paragraphs as paragraph, paragraphIndex (`${paragraphIndex}-${paragraph[0]?.verseNumber ?? 0}`)}
    <p>
      {#each paragraph as segment, index (`${paragraphIndex}-${segment.verseNumber}-${index}`)}
        {#if segment.showVerseNumber}
          <sup class="text-primary mr-0.5 text-xs">{segment.verseNumber}</sup>
        {/if}<FormattedVerseText text={segment.text} />
      {/each}
    </p>
  {/each}
</div>
