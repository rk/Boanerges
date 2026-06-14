<script lang="ts">
    import { bible } from '@/lib/bible.svelte.ts';

    let {
        label,
        value,
        onchange,
    }: {
        label: string;
        value: string;
        onchange: (id: string) => void;
    } = $props();
</script>

<label class="form-control w-full">
    <div class="label py-1">
        <span class="label-text text-xs">{label}</span>
    </div>
    <select
        class="select select-bordered select-sm w-full"
        {value}
        onchange={(event) => onchange(event.currentTarget.value)}
        disabled={bible.translationsLoading}
    >
        {#each bible.translations as translation (translation.id)}
            <option value={translation.id}>{translation.name}</option>
        {/each}
    </select>
</label>
