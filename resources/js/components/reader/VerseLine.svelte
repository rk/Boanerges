<script lang="ts">
    import Pilcrow from '@lucide/svelte/icons/pilcrow';

    import FormattedVerseText from '@/components/reader/FormattedVerseText.svelte';
    import type { Verse } from '@/lib/types/bible';

    let { verse, highlighted = false }: { verse: Verse; highlighted?: boolean } = $props();
</script>

<p class:mt-4={verse.paragraphStart} data-verse={verse.number}>
    {#if verse.paragraphStart}
        <Pilcrow size={14} class="mr-1 inline opacity-70" aria-hidden="true" />
    {/if}
    <span
        id="verse-{verse.number}"
        class={highlighted ? 'rounded-sm bg-primary/15 px-0.5 [box-decoration-break:clone]' : undefined}
    ><sup class="text-primary mr-1 text-xs">{verse.number}</sup><FormattedVerseText text={verse.text} /></span>
</p>
