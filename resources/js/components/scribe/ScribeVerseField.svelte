<script lang="ts">
    import type { Verse } from '@/lib/types/bible';

    let {
        verse,
        value = '',
        onupdate,
    }: {
        verse: Verse;
        value?: string;
        onupdate?: (value: string) => void;
    } = $props();

    const isEmpty = $derived(value.trim().length === 0);
</script>

<div class="mb-1" class:mt-4={verse.paragraphStart}>
    <div class="flex items-start gap-2">
        <span class="reader-prose mt-2 w-6 shrink-0 text-xs opacity-70">
            {#if verse.paragraphStart}¶{:else}{verse.number}{/if}
        </span>
        <textarea
            class="textarea reader-prose textarea-ghost w-full resize-none leading-relaxed"
            class:bg-base-300={isEmpty}
            rows="2"
            placeholder="Write your understanding of this verse…"
            {value}
            oninput={(event) => onupdate?.(event.currentTarget.value)}
        ></textarea>
    </div>
</div>
