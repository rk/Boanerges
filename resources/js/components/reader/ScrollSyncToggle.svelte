<script lang="ts">
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
    <span class="label-text text-sm">Link scroll</span>
    <input type="checkbox" class="toggle toggle-sm" checked={study.scrollSync} onchange={handleChange} />
</label>
