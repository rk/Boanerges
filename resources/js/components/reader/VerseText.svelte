<script lang="ts">
    import VerseLine from '@/components/reader/VerseLine.svelte';
    import type { Verse } from '@/lib/types/bible';

    let { verses, highlightedVerses = new Set<number>() }: {
        verses: Verse[];
        highlightedVerses?: Set<number>;
    } = $props();
</script>

<div class="reader-prose text-base-content">
    {#each verses as verse (verse.number)}
        <VerseLine {verse} highlighted={highlightedVerses.has(verse.number)} />
    {/each}
</div>
