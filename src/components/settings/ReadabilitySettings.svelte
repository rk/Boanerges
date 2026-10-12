<script lang="ts">
  import { readability, saveReadability } from "../../lib/readability.svelte.ts";

  async function changed() {
    await saveReadability();
  }
</script>

<div class="grid gap-3 p-4">
  <label class="form-control">
    <span class="label-text">Font size</span>
    <input class="range range-primary" type="range" min="14" max="28" bind:value={readability.fontSize} onchange={changed} />
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
</div>
