<script lang="ts">
  import { onMount } from "svelte";
  import { AppChrome, events, type PrintDocument } from "../mygo";

  let includeUserWork = $state(true);
  let busy = $state(false);
  let error = $state<string | null>(null);
  let savedPdfPath = $state<string | null>(null);
  let doc = $state<PrintDocument | null>(null);

  async function loadDocument(): Promise<void> {
    error = null;
    savedPdfPath = null;
    try {
      doc = await AppChrome.printPreviewDocument(includeUserWork);
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
      doc = null;
    }
  }

  async function runPrint(): Promise<void> {
    busy = true;
    error = null;
    try {
      await AppChrome.printPreviewPrint();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = false;
    }
  }

  async function runSavePdf(): Promise<void> {
    busy = true;
    error = null;
    savedPdfPath = null;
    try {
      savedPdfPath = await AppChrome.printPreviewSavePDF();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    } finally {
      busy = false;
    }
  }

  $effect(() => {
    includeUserWork;
    void loadDocument();
  });

  onMount(() => {
    const off = events.studyPrintPreviewRefresh.on(() => void loadDocument());
    return off;
  });
</script>

{#if doc}
  {@html `<style>${doc.styles}</style>`}
{/if}

<div
  class="box-border flex h-dvh flex-col gap-3 overflow-hidden bg-base-200 p-4 print:h-auto print:overflow-visible print:p-0"
  data-testid="print-preview"
>
  <div class="border-base-300 min-h-0 flex-1 overflow-y-auto rounded-lg border bg-white print:rounded-none print:border-0">
    {#if error}
      <p class="text-error p-4 text-sm">{error}</p>
    {:else if doc}
      {@html doc.body}
    {:else}
      <p class="p-4 text-sm opacity-60">Loading print preview…</p>
    {/if}
  </div>

  <div
    role="contentinfo"
    class="border-base-300 bg-base-100 box-border shrink-0 rounded-lg border print:hidden"
    data-testid="print-preview-footer"
  >
    <div class="box-border flex w-full flex-wrap items-center gap-x-4 gap-y-3 px-5 py-4">
    <label class="label cursor-pointer gap-3 !p-0">
      <input
        type="checkbox"
        class="toggle toggle-primary"
        bind:checked={includeUserWork}
        disabled={busy}
        data-testid="print-include-work"
      />
      <span class="label-text">
        <span class="block font-medium">Include my notes</span>
        <span class="block text-xs opacity-70">Off = lined notes and scribe for handwriting</span>
      </span>
    </label>
    <div class="flex-1"></div>
    {#if savedPdfPath}
      <p class="text-success text-xs">Saved {savedPdfPath}</p>
    {/if}
    <button type="button" class="btn btn-ghost" disabled={busy || !doc} data-testid="print-save-pdf" onclick={() => runSavePdf()}>
      Print to PDF
    </button>
    <button type="button" class="btn btn-primary" disabled={busy || !doc} data-testid="print-run" onclick={() => runPrint()}>
      Print
    </button>
    </div>
  </div>
</div>
