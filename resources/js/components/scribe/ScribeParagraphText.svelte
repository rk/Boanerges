<script lang="ts">
    import FormattedVerseText from '@/components/reader/FormattedVerseText.svelte';
    import { buildScribePreviewParagraphs } from '@/lib/paragraphs.ts';
    import type { Verse } from '@/lib/types/bible';

    let { verses }: { verses: Verse[] } = $props();

    const paragraphs = $derived(buildScribePreviewParagraphs(verses));
</script>

<div class="reader-prose text-base-content space-y-4">
    {#each paragraphs as paragraph, paragraphIndex (`${paragraphIndex}-${paragraph[0]?.verseNumber ?? 0}`)}
        <p>
            {#each paragraph as segment, index (`${paragraphIndex}-${segment.verseNumber}-${index}`)}
                {#if index > 0}{/if}
                {#if segment.showVerseNumber}
                    <sup class="text-primary mr-0.5 text-xs"
                        >{segment.verseNumber}</sup
                    >
                {/if}<FormattedVerseText text={segment.text} />
            {/each}
        </p>
    {/each}
</div>
