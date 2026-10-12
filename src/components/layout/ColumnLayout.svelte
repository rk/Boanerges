<script lang="ts">
  import { study } from "../../lib/study.svelte.ts";
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
</script>

<div class="flex h-[calc(100vh-3rem)] min-h-0">
  <section class="min-w-0 flex-1 border-r border-base-300">
    <BibleColumn translationId={study.translationId} primary />
  </section>
  {#each study.columns.slice(0, study.columnCount - 1) as column, i}
    {@const Component = slotComponents[column as keyof typeof slotComponents] ?? NotesColumn}
    <section class="min-w-0 flex-1 border-r border-base-300 last:border-r-0">
      <Component translationId={column === "bible-secondary" ? study.translationBId : study.translationId} slotIndex={i} columnType={column} />
    </section>
  {/each}
</div>
