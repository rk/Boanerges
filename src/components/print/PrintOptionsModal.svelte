<script lang="ts">
  import FileText from "@lucide/svelte/icons/file-text";
  import NotebookPen from "@lucide/svelte/icons/notebook-pen";
  import { printStudy, saveStudyPrintHtml, type PrintMode } from "../../lib/printStudy";
  import { syncDialog } from "../../lib/syncDialog";

  let { open = $bindable(false) }: { open?: boolean } = $props();
  let busy = $state(false);
  let error = $state<string | null>(null);
  let savedPath = $state<string | null>(null);

  function onPresent(): void {
    error = null;
    savedPath = null;
  }

  async function saveHtml(includeUserWork: boolean): Promise<void> {
    busy = true;
    error = null;
    try {
      savedPath = await saveStudyPrintHtml(includeUserWork);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = false;
    }
  }

  async function runPrint(mode: PrintMode): Promise<void> {
    busy = true;
    error = null;
    try {
      await printStudy(mode);
      open = false;
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = false;
    }
  }

  function onDialogClose(): void {
    open = false;
  }
</script>

<dialog
  use:syncDialog={{ open, onPresent }}
  class="modal"
  data-testid="print-options"
  onclose={onDialogClose}
>
  <div class="modal-box max-w-md">
    <h3 class="font-bold text-lg">Print chapter</h3>
    <p class="text-sm opacity-70 mt-1">
      Bible columns print the chapter text. Scribe always prints as lined space for handwriting.
    </p>
    {#if error}
      <p class="text-error text-sm mt-2">{error}</p>
    {/if}
    {#if savedPath}
      <p class="text-success text-sm mt-2">Saved to {savedPath}</p>
    {/if}
    <div class="mt-4 flex flex-col gap-3">
      <button
        type="button"
        class="btn justify-start gap-3"
        disabled={busy}
        data-testid="print-include-work"
        onclick={() => runPrint("include-user-work")}
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
        onclick={() => runPrint("blank-writing")}
      >
        <FileText size={18} aria-hidden="true" />
        <span class="text-left">
          <span class="block font-medium">Blank for writing</span>
          <span class="block text-xs opacity-70">Lined notes and scribe areas only</span>
        </span>
      </button>
    </div>
    <div class="modal-action flex-wrap gap-2">
      <button type="button" class="btn btn-ghost btn-sm" disabled={busy} data-testid="print-save-html" onclick={() => saveHtml(true)}>
        Save HTML…
      </button>
      <button type="button" class="btn btn-ghost" disabled={busy} onclick={() => (open = false)}>Cancel</button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button type="submit">close</button></form>
</dialog>
