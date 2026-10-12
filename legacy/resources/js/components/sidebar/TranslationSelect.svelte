<script lang="ts">
    import { bible } from '@/lib/bible.svelte.ts';

    let {
        value,
        onchange,
    }: {
        value: string;
        onchange: (id: string) => void;
    } = $props();
</script>

<select
    class="select select-bordered select-sm grow"
    {value}
    onchange={(event) => onchange(event.currentTarget.value)}
    disabled={bible.translationsLoading}
    aria-label="Translation"
>
    {#each bible.translations as translation (translation.id)}
        <option value={translation.id}>{translation.name}</option>
    {/each}
</select>
