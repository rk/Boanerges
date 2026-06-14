<script lang="ts">
    import VerseLine from '@/components/reader/VerseLine.svelte';
    import type { Verse } from '@/lib/types/bible';

    let { verses }: { verses: Verse[] } = $props();
</script>

<div class="reader-prose text-base-content">
    {#each verses as verse (verse.number)}
        <VerseLine {verse} />
    {/each}
</div>
