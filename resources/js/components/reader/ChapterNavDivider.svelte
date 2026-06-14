<script lang="ts">
    import ChevronDown from '@lucide/svelte/icons/chevron-down';
    import ChevronUp from '@lucide/svelte/icons/chevron-up';

    let {
        direction,
        bookAbbrev,
        chapter,
        layout = 'horizontal',
        onclick,
    }: {
        direction: 'prev' | 'next';
        bookAbbrev: string;
        chapter: number;
        layout?: 'horizontal' | 'vertical';
        onclick?: () => void;
    } = $props();

    const Icon = $derived(direction === 'prev' ? ChevronUp : ChevronDown);
</script>

{#if onclick}
    <button
        type="button"
        class="divider hover:text-primary flex cursor-pointer items-center gap-1 text-sm opacity-80"
        class:divider-horizontal={layout === 'horizontal'}
        class:divider-vertical={layout === 'vertical'}
        class:flex-col={layout === 'horizontal'}
        {onclick}
    >
        <Icon color="currentColor" size={18} aria-hidden="true" />
        <span>{bookAbbrev} {chapter}</span>
        <Icon color="currentColor" size={18} aria-hidden="true" />
    </button>
{:else}
    <div
        class="divider flex items-center gap-1 text-sm opacity-50"
        class:divider-horizontal={layout === 'horizontal'}
        class:divider-vertical={layout === 'vertical'}
        class:flex-col={layout === 'horizontal'}
    >
        <Icon color="currentColor" size={18} aria-hidden="true" />
        <span>{bookAbbrev} {chapter}</span>
        <Icon color="currentColor" size={18} aria-hidden="true" />
    </div>
{/if}
