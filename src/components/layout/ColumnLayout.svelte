<script lang="ts">
  import { syncScrollFromSource } from "../../lib/scrollSync.ts";
  import { study } from "../../lib/study.svelte.ts";
  import { secondaryTranslationForSlot } from "../../lib/studyLayout";
  import BibleColumn from "../columns/BibleColumn.svelte";
  import ComparisonColumn from "../columns/ComparisonColumn.svelte";
  import CrossReferencesColumn from "../columns/CrossReferencesColumn.svelte";
  import DictionaryColumn from "../columns/DictionaryColumn.svelte";
  import NotesColumn from "../columns/NotesColumn.svelte";
  import ScribeColumn from "../columns/ScribeColumn.svelte";
  import SearchColumn from "../columns/SearchColumn.svelte";
  import VerseListColumn from "../columns/VerseListColumn.svelte";

  const slotComponents = {
    notes: NotesColumn,
    scribe: ScribeColumn,
    search: SearchColumn,
    "cross-references": CrossReferencesColumn,
    comparison: ComparisonColumn,
    "verse-list": VerseListColumn,
    dictionary: DictionaryColumn,
    "bible-secondary": BibleColumn,
  } as const;

  let primaryScroll = $state<HTMLElement | null>(null);
  let secondaryScroll = $state<HTMLElement | null>(null);

  function onPrimaryScroll() {
    if (!study.scrollSync || !primaryScroll) return;
    const targets: HTMLElement[] = [];
    if (secondaryScroll) targets.push(secondaryScroll);
    syncScrollFromSource(primaryScroll, targets);
  }

  function secondaryTranslationId(slotIndex: number): string {
    return secondaryTranslationForSlot(slotIndex, study.translationBId, study.translationCId);
  }
</script>

<div class="flex h-[calc(100vh-3rem)] min-h-0" data-testid="column-layout">
  <section class="min-w-0 flex-1 border-r border-base-300">
    <BibleColumn translationId={study.translationId} primary bind:scrollRef={primaryScroll} onscroll={onPrimaryScroll} />
  </section>
  {#each study.columns.slice(0, study.columnCount - 1) as column, i}
    {@const Component = slotComponents[column as keyof typeof slotComponents] ?? NotesColumn}
    <section class="min-w-0 flex-1 border-r border-base-300 last:border-r-0">
      {#if column === "bible-secondary"}
        <BibleColumn
          translationId={secondaryTranslationId(i)}
          slotIndex={i}
          columnType={column}
          bind:scrollRef={secondaryScroll}
          onscroll={onPrimaryScroll}
        />
      {:else if column === "comparison"}
        <ComparisonColumn slotIndex={i} />
      {:else}
        <Component slotIndex={i} columnType={column} />
      {/if}
    </section>
  {/each}
</div>
