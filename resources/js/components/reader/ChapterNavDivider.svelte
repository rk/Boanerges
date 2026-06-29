<script lang="ts">
    import { ArrowUp, ArrowDown } from '@lucide/svelte';

    let {
        direction,
        label,
        layout = 'compact',
        onclick,
    }: {
        direction: 'prev' | 'next';
        label: string;
        layout?: 'compact' | 'divider';
        onclick?: () => void;
    } = $props();

    const Icon = $derived(direction === 'prev' ? ArrowUp : ArrowDown);
</script>

{#if onclick}
    <button
        type="button"
        class="hover:text-primary flex cursor-pointer items-center gap-1 text-sm opacity-80"
        class:divider={layout === 'divider'}
        class:divider-vertical={layout === 'divider'}
        class:flex-col={layout === 'compact'}
        {onclick}
    >
        {#if layout === 'compact'}
            {#if direction === 'next'}{label}{/if}
            <Icon color="currentColor" size={16} aria-hidden="true" />
            {#if direction === 'prev'}{label}{/if}
        {:else}
            <Icon color="currentColor" size={18} aria-hidden="true" />
            <span>{label}</span>
            <Icon color="currentColor" size={18} aria-hidden="true" />
        {/if}
    </button>
{:else}
    <div
        class="divider flex items-center gap-1 text-sm opacity-50"
        class:divider={layout === 'divider'}
        class:divider-vertical={layout === 'divider'}
        class:flex-col={layout === 'compact'}
    >
        {#if layout === 'compact'}
            {#if direction === 'next'}{label}{/if}
            <Icon color="currentColor" size={16} aria-hidden="true" />
            {#if direction === 'prev'}{label}{/if}
        {:else}
            <Icon color="currentColor" size={18} aria-hidden="true" />
            <span>{label}</span>
            <Icon color="currentColor" size={18} aria-hidden="true" />
        {/if}
    </div>
{/if}
