<script lang="ts">
    import Pilcrow from '@lucide/svelte/icons/pilcrow';

    import type { Verse } from '@/lib/types/bible';

    let {
        verse,
        value = '',
        paragraphStart = false,
        onupdate,
        onToggleParagraph,
    }: {
        verse: Verse;
        value?: string;
        paragraphStart?: boolean;
        onupdate?: (value: string) => void;
        onToggleParagraph?: () => void;
    } = $props();

    const isEmpty = $derived(value.trim().length === 0);
</script>

<div class="mb-1" class:mt-4={paragraphStart}>
    <div class="flex items-start gap-2">
        <span class="reader-prose mt-2 flex w-6 shrink-0 flex-col items-center gap-1 text-xs opacity-70">
            <button
                type="button"
                class="btn btn-ghost btn-xs px-0"
                class:opacity-30={! paragraphStart}
                class:opacity-100={paragraphStart}
                disabled={verse.number === 1}
                aria-label={paragraphStart ? 'Remove paragraph break before this verse' : 'Start a new paragraph before this verse'}
                onclick={() => onToggleParagraph?.()}
            >
                <Pilcrow size={14} aria-hidden="true" />
            </button>
            {verse.number}
        </span>
        <textarea
            class="textarea reader-prose textarea-ghost w-full resize-y leading-relaxed"
            class:bg-base-300={isEmpty}
            rows="2"
            placeholder="Write your revision of this verse. Use a blank line for a new paragraph within this verse."
            {value}
            oninput={(event) => onupdate?.(event.currentTarget.value)}
        ></textarea>
    </div>
</div>
