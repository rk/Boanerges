<script lang="ts">
    import Printer from '@lucide/svelte/icons/printer';

    import PrintOptionsModal from '@/components/print/PrintOptionsModal.svelte';
    import { openPrintOptions, printUi } from '@/lib/printStudy.svelte.ts';
</script>

<div class="fab fixed right-4 bottom-4 z-40 print:hidden">
    <button
        type="button"
        class="btn btn-lg btn-circle btn-primary shadow-lg"
        aria-label="Print chapter"
        onclick={openPrintOptions}
    >
        <Printer size={22} aria-hidden="true" />
    </button>
</div>

<PrintOptionsModal
    open={printUi.optionsOpen}
    onclose={() => (printUi.optionsOpen = false)}
/>
