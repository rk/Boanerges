<script lang="ts">
    import type { Verse } from '@/lib/types/bible';

    let { verses }: { verses: Verse[] } = $props();

    const paragraphs = $derived.by(() => {
        const groups: Verse[][] = [];
        let current: Verse[] = [];

        for (const verse of verses) {
            if (verse.paragraphStart && current.length > 0) {
                groups.push(current);
                current = [];
            }

            current.push(verse);
        }

        if (current.length > 0) {
            groups.push(current);
        }

        return groups;
    });
</script>

<div class="reader-prose text-base-content space-y-4">
    {#each paragraphs as paragraph (paragraph[0].number)}
        <p>
            {#each paragraph as verse, index (verse.number)}
                {#if index > 0}
                    {' '}
                {/if}
                <sup class="text-primary mr-0.5 text-xs">{verse.number}</sup>{verse.text}
            {/each}
        </p>
    {/each}
</div>
