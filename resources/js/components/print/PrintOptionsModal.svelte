<script lang="ts">
    import FileText from '@lucide/svelte/icons/file-text';
    import NotebookPen from '@lucide/svelte/icons/notebook-pen';

    import {
        fetchStudyPrinters,
        printStudy,
        PRINT_TO_PDF,
    } from '@/lib/printStudy.ts';
    import type { PrintMode, StudyPrinter } from '@/lib/printStudy.ts';

    let {
        open = false,
        onclose,
    }: {
        open?: boolean;
        onclose?: () => void;
    } = $props();

    let dialog: HTMLDialogElement | undefined = $state();
    let printing = $state(false);
    let loadingPrinters = $state(false);
    let error = $state<string | null>(null);
    let success = $state<string | null>(null);
    let printers = $state<StudyPrinter[]>([]);
    let selectedPrinterName = $state('');

    $effect(() => {
        if (open) {
            error = null;
            success = null;
            selectedPrinterName = '';
            dialog?.showModal();
            void loadPrinters();
        } else {
            dialog?.close();
        }
    });

    async function loadPrinters(): Promise<void> {
        loadingPrinters = true;

        try {
            printers = await fetchStudyPrinters();
        } catch (caught) {
            printers = [];
            error =
                caught instanceof Error
                    ? caught.message
                    : 'Could not load printers.';
        } finally {
            loadingPrinters = false;
        }
    }

    function handleClose(): void {
        if (printing) {
            return;
        }

        onclose?.();
    }

    async function handlePrint(mode: PrintMode): Promise<void> {
        printing = true;
        error = null;
        success = null;

        try {
            const path = await printStudy(mode, selectedPrinterName);

            if (path) {
                success = `PDF saved to ${path}`;

                return;
            }

            onclose?.();
        } catch (caught) {
            error = caught instanceof Error ? caught.message : 'Print failed.';
        } finally {
            printing = false;
        }
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box max-w-md">
        <h2 class="text-lg font-semibold">Print chapter</h2>
        <p class="text-base-content/70 mt-2 text-sm">
            Bible columns print the chapter text. Scribe always prints as lined
            space for handwriting.
        </p>

        <fieldset class="fieldset mt-4">
            <legend class="fieldset-legend">Printer</legend>
            {#if loadingPrinters}
                <span class="loading loading-spinner loading-sm text-primary"
                ></span>
            {:else}
                <select
                    class="select select-bordered w-full"
                    bind:value={selectedPrinterName}
                    disabled={printing}
                    aria-label="Printer"
                >
                    <option value="">System default</option>
                    <option value={PRINT_TO_PDF}>Print to PDF</option>
                    {#each printers as printer (printer.name)}
                        <option value={printer.name}
                            >{printer.displayName}</option
                        >
                    {/each}
                </select>
            {/if}
        </fieldset>

        {#if success}
            <div class="alert alert-success mt-4 text-sm">
                <span>{success}</span>
            </div>
        {/if}

        {#if error}
            <div class="alert alert-error mt-4 text-sm">
                <span>{error}</span>
            </div>
        {/if}

        <div class="mt-6 flex flex-col gap-3">
            <button
                type="button"
                class="btn justify-start gap-3"
                disabled={printing || loadingPrinters}
                onclick={() => handlePrint('include-user-work')}
            >
                <NotebookPen size={18} aria-hidden="true" />
                <span class="text-left">
                    <span class="block font-medium">Include my notes</span>
                    <span class="text-base-content/70 block text-xs"
                        >Notes content included; scribe stays blank</span
                    >
                </span>
            </button>

            <button
                type="button"
                class="btn btn-outline justify-start gap-3"
                disabled={printing || loadingPrinters}
                onclick={() => handlePrint('blank-writing')}
            >
                <FileText size={18} aria-hidden="true" />
                <span class="text-left">
                    <span class="block font-medium">Blank for writing</span>
                    <span class="text-base-content/70 block text-xs"
                        >Lined notes and scribe areas only</span
                    >
                </span>
            </button>
        </div>

        <div class="modal-action">
            <form method="dialog">
                <button type="submit" class="btn" disabled={printing}
                    >Cancel</button
                >
            </form>
        </div>
    </div>

    <form method="dialog" class="modal-backdrop">
        <button
            type="submit"
            aria-label="Close print options"
            disabled={printing}>Close</button
        >
    </form>
</dialog>
