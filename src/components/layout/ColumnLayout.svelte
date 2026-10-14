<script lang="ts">
  import { syncLinkedScrollPanes } from "../../lib/scrollSync.ts";
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
  let slotScrollRefs = $state<Record<number, HTMLElement | null>>({});
  let syncingScroll = false;

  function setSlotScrollRef(index: number, el: HTMLElement | null) {
    if (slotScrollRefs[index] === el) {
      return;
    }
    slotScrollRefs = { ...slotScrollRefs, [index]: el };
  }

  function bibleScrollPanes(): HTMLElement[] {
    const panes: HTMLElement[] = [];
    if (primaryScroll) {
      panes.push(primaryScroll);
    }
    const secondaryColumns = study.columns.slice(0, Math.max(0, study.columnCount - 1));
    for (let i = 0; i < secondaryColumns.length; i++) {
      if (secondaryColumns[i] === "bible-secondary") {
        const el = slotScrollRefs[i];
        if (el) {
          panes.push(el);
        }
      }
    }
    return panes;
  }

  function onLinkedScroll(source: HTMLElement) {
    if (!study.scrollSync || syncingScroll) {
      return;
    }
    const panes = bibleScrollPanes();
    if (panes.length < 2) {
      return;
    }
    syncingScroll = true;
    try {
      syncLinkedScrollPanes(source, panes);
    } finally {
      syncingScroll = false;
    }
  }

  function secondaryTranslationId(slotIndex: number): string {
    return secondaryTranslationForSlot(slotIndex, study.translationBId, study.translationCId);
  }
</script>

<div class="flex h-full min-h-0" data-testid="column-layout">
  <section class="min-w-0 flex-1 border-r border-base-300">
    <BibleColumn translationId={study.translationId} primary bind:scrollRef={primaryScroll} onscroll={onLinkedScroll} />
  </section>
  {#each study.columns.slice(0, study.columnCount - 1) as column, i}
    {@const Component = slotComponents[column as keyof typeof slotComponents] ?? NotesColumn}
    <section class="min-w-0 flex-1 border-r border-base-300 last:border-r-0">
      {#if column === "bible-secondary"}
        <BibleColumn
          translationId={secondaryTranslationId(i)}
          slotIndex={i}
          columnType={column}
          onScrollElement={(el) => setSlotScrollRef(i, el)}
          onscroll={onLinkedScroll}
        />
      {:else if column === "comparison"}
        <ComparisonColumn slotIndex={i} />
      {:else}
        <Component slotIndex={i} columnType={column} />
      {/if}
    </section>
  {/each}
</div>
