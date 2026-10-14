<script lang="ts">
  import Pilcrow from "@lucide/svelte/icons/pilcrow";
  import type { Verse } from "../../mygo";
  import VerseSpan from "./VerseSpan.svelte";

  let {
    verse,
    bookId,
    chapter,
    showContextMenu = false,
  }: {
    verse: Verse;
    bookId: string;
    chapter: number;
    showContextMenu?: boolean;
  } = $props();
</script>

<p class:mt-4={verse.paragraphStart} data-verse={verse.number}>
  {#if verse.paragraphStart}
    <Pilcrow size={14} class="mr-1 inline opacity-70" aria-hidden="true" />
  {/if}
  <VerseSpan {verse} {bookId} {chapter} {showContextMenu} />
</p>
