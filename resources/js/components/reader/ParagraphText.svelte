<script lang="ts">
    import FormattedVerseText from '@/components/reader/FormattedVerseText.svelte';
    import { groupVersesIntoParagraphs } from '@/lib/paragraphs.ts';
    import type { Verse } from '@/lib/types/bible';

    let { verses, highlightedVerses = new Set<number>() }: {
        verses: Verse[];
        highlightedVerses?: Set<number>;
    } = $props();

    const paragraphs = $derived(groupVersesIntoParagraphs(verses));
</script>

<div class="reader-prose text-base-content space-y-4">
    {#each paragraphs as paragraph (paragraph[0].number)}
        <p>
            {#each paragraph as verse, index (verse.number)}
                {#if index > 0}
                    {' '}
                {/if}
                <span
                    id="verse-{verse.number}"
                    data-verse={verse.number}
                    class={highlightedVerses.has(verse.number) ? 'rounded-sm bg-primary/15 px-0.5 [box-decoration-break:clone]' : undefined}
                ><sup class="text-primary mr-0.5 text-xs">{verse.number}</sup><FormattedVerseText text={verse.text} /></span>
            {/each}
        </p>
    {/each}
</div>
