<script lang="ts">
    import Pilcrow from '@lucide/svelte/icons/pilcrow';

    import FormattedVerseText from '@/components/reader/FormattedVerseText.svelte';
    import type { Verse } from '@/lib/types/bible';

    let { verse }: { verse: Verse } = $props();
</script>

<p class:mt-4={verse.paragraphStart}>
    {#if verse.paragraphStart}
        <Pilcrow size={14} class="mr-1 inline opacity-70" aria-hidden="true" />
    {/if}
    <sup class="text-primary mr-1 text-xs">{verse.number}</sup><FormattedVerseText text={verse.text} />
</p>
