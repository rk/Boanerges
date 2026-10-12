<script lang="ts">
  import { parseVerseHtml, type VerseHtmlNode } from "../../lib/parseVerseHtml";

  let { text }: { text: string } = $props();

  const nodes = $derived(parseVerseHtml(text));
</script>

{#each nodes as node, index (index)}
  {@render renderNode(node)}
{/each}

{#snippet renderNode(node: VerseHtmlNode)}
  {#if node.type === "text"}
    {node.value}
  {:else if node.name === "em"}
    <em>{#each node.children as child, i (i)}{@render renderNode(child)}{/each}</em>
  {:else if node.name === "strong"}
    <strong>{#each node.children as child, i (i)}{@render renderNode(child)}{/each}</strong>
  {:else if node.name === "sup"}
    <sup>{#each node.children as child, i (i)}{@render renderNode(child)}{/each}</sup>
  {:else}
    <span class={node.className}>{#each node.children as child, i (i)}{@render renderNode(child)}{/each}</span>
  {/if}
{/snippet}
