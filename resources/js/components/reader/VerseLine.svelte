<script lang="ts">
    import type { Verse } from '@/lib/types/bible';

    let { verse }: { verse: Verse } = $props();
</script>

<p class:mt-4={verse.paragraphStart}>
    {#if verse.paragraphStart}
        <span class="mr-1 opacity-70">¶</span>
    {/if}
    <sup class="text-primary mr-1 text-xs">{verse.number}</sup>{verse.text}
</p>
