<script lang="ts">
    let {
        label,
        direction = 'horizontal',
        onclick,
    }: {
        label: string;
        direction?: 'horizontal' | 'vertical';
        onclick?: () => void;
    } = $props();
</script>

{#if onclick}
    <button
        type="button"
        class="divider hover:text-primary cursor-pointer text-sm opacity-80"
        class:divider-horizontal={direction === 'horizontal'}
        class:divider-vertical={direction === 'vertical'}
        {onclick}
    >
        {label}
    </button>
{:else}
    <div
        class="divider text-sm opacity-50"
        class:divider-horizontal={direction === 'horizontal'}
        class:divider-vertical={direction === 'vertical'}
    >
        {label}
    </div>
{/if}
