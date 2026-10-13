<script lang="ts">
  import { onMount } from "svelte";
  import { AppChrome, Service } from "../mygo";
  import BookChapterPickerPanel from "../components/toolbar/BookChapterPickerPanel.svelte";
  import { hydrateStudy, schedulePersistStudy, study, type ColumnType } from "../lib/study.svelte.ts";
  import { normalizeColumns } from "../lib/studyLayout";
  import { hydrateReadability } from "../lib/readability.svelte.ts";

  let step = $state<"book" | "chapter">("book");

  onMount(() => {
    void Service.getStudy().then(hydrateStudy);
    void Service.getReadability().then(hydrateReadability);
  });

  async function onDone(): Promise<void> {
    schedulePersistStudy();
    const columns = normalizeColumns(study.columnCount, study.columns as ColumnType[]);
    await Service.updateStudy({
      columnCount: study.columnCount,
      columns,
      bookId: study.bookId,
      chapter: study.chapter,
      translationId: study.translationId,
      translationBId: study.translationBId,
      translationCId: study.translationCId,
      verseListShowContent: study.verseListShowContent,
      verseListActiveId: study.verseListActiveId,
    });
    await AppChrome.closeBookChapterFlyout();
  }
</script>

<BookChapterPickerPanel bind:step {onDone} />
