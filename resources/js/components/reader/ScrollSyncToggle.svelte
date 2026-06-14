<script lang="ts">
    import Link2 from '@lucide/svelte/icons/link-2';
    import Unlink from '@lucide/svelte/icons/unlink';

    import { study, setScrollSync } from '@/lib/study.svelte.ts';

    let {
        onchange,
    }: {
        onchange?: (enabled: boolean) => void;
    } = $props();

    function handleChange(event: Event): void {
        const enabled = (event.currentTarget as HTMLInputElement).checked;
        setScrollSync(enabled);
        onchange?.(enabled);
    }
</script>

<label class="label cursor-pointer gap-2">
    {#if study.scrollSync}
        <Link2 size={16} aria-hidden="true" />
    {:else}
        <Unlink size={16} aria-hidden="true" />
    {/if}
    <span class="label-text text-sm">Link scroll</span>
    <input type="checkbox" class="toggle toggle-sm" checked={study.scrollSync} onchange={handleChange} />
</label>
