<script lang="ts">
    import { ArrowUp, ArrowDown } from '@lucide/svelte';

    let {
        direction,
        label,
        layout = 'horizontal',
        onclick
    }: {
        direction: 'prev' | 'next';
        label: string;
        layout?: 'horizontal' | 'vertical';
        onclick?: () => void;
    } = $props();

    const Icon = $derived(direction === 'prev' ? ArrowUp : ArrowDown);
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
        <span>{label}</span>
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
        <span>{label}</span>
        <Icon color="currentColor" size={18} aria-hidden="true" />
    </div>
{/if}
