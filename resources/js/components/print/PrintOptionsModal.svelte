<script lang="ts">
    import FileText from '@lucide/svelte/icons/file-text';
    import NotebookPen from '@lucide/svelte/icons/notebook-pen';

    import { printStudy, type PrintMode } from '@/lib/printStudy.ts';

    let {
        open = false,
        onclose,
    }: {
        open?: boolean;
        onclose?: () => void;
    } = $props();

    let dialog: HTMLDialogElement | undefined = $state();
    let printing = $state(false);
    let error = $state<string | null>(null);

    $effect(() => {
        if (open) {
            error = null;
            dialog?.showModal();
        } else {
            dialog?.close();
        }
    });

    function handleClose(): void {
        if (printing) {
            return;
        }

        onclose?.();
    }

    async function handlePrint(mode: PrintMode): Promise<void> {
        printing = true;
        error = null;

        try {
            await printStudy(mode);
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
            Bible columns print the chapter text. Scribe always prints as lined space for handwriting.
        </p>

        {#if error}
            <div class="alert alert-error mt-4 text-sm">
                <span>{error}</span>
            </div>
        {/if}

        <div class="mt-6 flex flex-col gap-3">
            <button
                type="button"
                class="btn justify-start gap-3"
                disabled={printing}
                onclick={() => handlePrint('include-user-work')}
            >
                <NotebookPen size={18} aria-hidden="true" />
                <span class="text-left">
                    <span class="block font-medium">Include my notes</span>
                    <span class="text-base-content/70 block text-xs">Notes content included; scribe stays blank</span>
                </span>
            </button>

            <button
                type="button"
                class="btn btn-outline justify-start gap-3"
                disabled={printing}
                onclick={() => handlePrint('blank-writing')}
            >
                <FileText size={18} aria-hidden="true" />
                <span class="text-left">
                    <span class="block font-medium">Blank for writing</span>
                    <span class="text-base-content/70 block text-xs">Lined notes and scribe areas only</span>
                </span>
            </button>
        </div>

        <div class="modal-action">
            <form method="dialog">
                <button type="submit" class="btn" disabled={printing}>Cancel</button>
            </form>
        </div>
    </div>

    <form method="dialog" class="modal-backdrop">
        <button type="submit" aria-label="Close print options" disabled={printing}>Close</button>
    </form>
</dialog>
