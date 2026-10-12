<script lang="ts">
  import { readability, saveReadability } from "../../lib/readability.svelte.ts";

  let error = $state<string | null>(null);

  async function changed() {
    error = null;
    try {
      await saveReadability();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function simulateInvalid() {
    readability.fontSize = 10;
    await changed();
  }
</script>

<div class="grid gap-3 p-4" data-testid="readability-settings">
  {#if error}
    <p class="text-error text-sm" data-testid="readability-error">{error}</p>
  {/if}
  <label class="form-control">
    <span class="label-text">Font size</span>
    <input class="range range-primary" type="range" min="14" max="24" bind:value={readability.fontSize} onchange={changed} />
  </label>
  <label class="form-control">
    <span class="label-text">Line height</span>
    <input class="range range-primary" type="range" min="12" max="24" bind:value={readability.lineHeight} onchange={changed} />
  </label>
  <label class="form-control">
    <span class="label-text">Theme</span>
    <select class="select select-bordered" bind:value={readability.theme} onchange={changed}>
      <option value="auto">Auto</option>
      <option value="light">Light</option>
      <option value="dark">Dark</option>
    </select>
  </label>
  <label class="label cursor-pointer justify-start gap-3">
    <input type="checkbox" class="checkbox" bind:checked={readability.justifyText} onchange={changed} />
    <span class="label-text">Justify text</span>
  </label>
  {#if import.meta.env.VITE_E2E}
    <button type="button" class="btn btn-xs btn-outline" data-testid="readability-invalid" onclick={simulateInvalid}>Simulate invalid</button>
  {/if}
</div>
