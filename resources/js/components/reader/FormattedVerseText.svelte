<script lang="ts">
    import { parseVerseHtml } from '@/lib/parseVerseHtml';
    import type { VerseHtmlNode } from '@/lib/parseVerseHtml';
    import { splitVerseTextForDictionary } from '@/lib/splitVerseTextForDictionary';

    let { text }: { text: string } = $props();

    const nodes = $derived(parseVerseHtml(text));
</script>

{#each nodes as node (node)}
    {@render renderNode(node)}
{/each}

{#snippet renderNode(node: VerseHtmlNode)}
    {#if node.type === 'text'}
        {#each splitVerseTextForDictionary(node.value) as part, index (`${index}-${part.value}`)}
            {#if part.kind === 'word'}
                <span data-dict-word={part.value}>{part.value}</span>
            {:else}
                {part.value}
            {/if}
        {/each}
    {:else if node.name === 'em'}
        <em
            >{#each node.children as child (child)}{@render renderNode(
                    child,
                )}{/each}</em
        >
    {:else if node.name === 'strong'}
        <strong
            >{#each node.children as child (child)}{@render renderNode(
                    child,
                )}{/each}</strong
        >
    {:else if node.name === 'sup'}
        <sup
            >{#each node.children as child (child)}{@render renderNode(
                    child,
                )}{/each}</sup
        >
    {:else}
        <span class={node.className}
            >{#each node.children as child (child)}{@render renderNode(
                    child,
                )}{/each}</span
        >
    {/if}
{/snippet}
