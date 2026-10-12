<script lang="ts">
  import FileText from "@lucide/svelte/icons/file-text";
  import NotebookPen from "@lucide/svelte/icons/notebook-pen";
  import { AppChrome } from "../../mygo";
  import { study } from "../../lib/study.svelte.ts";
  import { readability } from "../../lib/readability.svelte.ts";
  import type { StudySettings } from "../../mygo";

  let { open = $bindable(false) }: { open?: boolean } = $props();
  let dialog = $state<HTMLDialogElement | null>(null);
  let busy = $state(false);
  let error = $state<string | null>(null);

  $effect(() => {
    if (open) {
      error = null;
      dialog?.showModal();
    } else {
      dialog?.close();
    }
  });

  function studyPayload(): StudySettings {
    return {
      columnCount: study.columnCount,
      columns: study.columns,
      bookId: study.bookId,
      chapter: study.chapter,
      translationId: study.translationId,
      translationBId: study.translationBId,
      translationCId: study.translationCId,
      verseListShowContent: study.verseListShowContent,
      verseListActiveId: study.verseListActiveId,
    };
  }

  async function printStudy(includeUserWork: boolean): Promise<void> {
    busy = true;
    error = null;
    try {
      const html = await AppChrome.studyPrintHTML(
        studyPayload(),
        {
          fontSize: readability.fontSize,
          lineHeight: readability.lineHeight,
          theme: readability.theme,
          fontFamily: readability.fontFamily,
          justifyText: readability.justifyText,
        },
        includeUserWork,
      );
      const printWindow = window.open("", "_blank", "noopener,noreferrer");
      if (!printWindow) {
        throw new Error("Could not open print window.");
      }
      printWindow.document.open();
      printWindow.document.write(html);
      printWindow.document.close();
      printWindow.focus();
      printWindow.print();
      dialog?.close();
      open = false;
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = false;
    }
  }
</script>

<dialog bind:this={dialog} class="modal" data-testid="print-options" onclose={() => (open = false)}>
  <div class="modal-box max-w-md">
    <h3 class="font-bold text-lg">Print chapter</h3>
    <p class="text-sm opacity-70 mt-1">
      Bible columns print the chapter text. Scribe always prints as lined space for handwriting.
    </p>
    {#if error}
      <p class="text-error text-sm mt-2">{error}</p>
    {/if}
    <div class="mt-4 flex flex-col gap-3">
      <button
        type="button"
        class="btn justify-start gap-3"
        disabled={busy}
        data-testid="print-include-work"
        onclick={() => printStudy(true)}
      >
        <NotebookPen size={18} aria-hidden="true" />
        <span class="text-left">
          <span class="block font-medium">Include my notes</span>
          <span class="block text-xs opacity-70">Notes and scribe content included for reading</span>
        </span>
      </button>
      <button
        type="button"
        class="btn btn-outline justify-start gap-3"
        disabled={busy}
        data-testid="print-blank-writing"
        onclick={() => printStudy(false)}
      >
        <FileText size={18} aria-hidden="true" />
        <span class="text-left">
          <span class="block font-medium">Blank for writing</span>
          <span class="block text-xs opacity-70">Lined notes and scribe areas only</span>
        </span>
      </button>
    </div>
    <div class="modal-action">
      <button type="button" class="btn btn-ghost" disabled={busy} onclick={() => (open = false)}>Cancel</button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button type="submit">close</button></form>
</dialog>
