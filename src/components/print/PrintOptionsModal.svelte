<script lang="ts">
  import { AppChrome } from "../../mygo";

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

  async function printStudy() {
    busy = true;
    error = null;
    try {
      await AppChrome.print();
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
  <div class="modal-box">
    <h3 class="font-bold text-lg">Print study</h3>
    <p class="text-sm opacity-70 mt-1">Print visible columns using the system print dialog.</p>
    {#if error}
      <p class="text-error text-sm mt-2">{error}</p>
    {/if}
    <div class="modal-action">
      <button type="button" class="btn btn-ghost" onclick={() => (open = false)}>Cancel</button>
      <button type="button" class="btn btn-primary" disabled={busy} onclick={printStudy}>Print…</button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button type="submit">close</button></form>
</dialog>
