<script lang="ts">
  import { DictionaryService } from "../../mygo";
  import { readerFontStack, type ReaderFontFamily } from "../../lib/readerFonts";
  import { readability, saveReadability, syncReaderTheme } from "../../lib/readability.svelte.ts";

  const fontOptions: { id: ReaderFontFamily; label: string }[] = [
    { id: "sans-serif", label: "Sans-serif" },
    { id: "serif", label: "Serif" },
  ];

  function setFontFamily(id: ReaderFontFamily) {
    readability.fontFamily = id;
    void changed();
  }

  let error = $state<string | null>(null);
  let dictionaryMsg = $state<string | null>(null);
  let dictionaryBusy = $state(false);

  async function changed() {
    error = null;
    syncReaderTheme();
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

  async function refreshDictionaryStatus() {
    try {
      const st = (await DictionaryService.status()) as { installed?: boolean; entries?: number };
      if (st.installed && (st.entries ?? 0) > 0) {
        dictionaryMsg = `Webster loaded (${st.entries?.toLocaleString() ?? "?"} entries).`;
      } else {
        dictionaryMsg = `Webster not ready (${st.entries ?? 0} entries). Use Reinstall or look up a word.`;
      }
    } catch (err) {
      dictionaryMsg = err instanceof Error ? err.message : String(err);
    }
  }

  async function reinstallDictionary() {
    dictionaryMsg = null;
    dictionaryBusy = true;
    try {
      await DictionaryService.reinstall();
      await refreshDictionaryStatus();
    } catch (err) {
      dictionaryMsg = err instanceof Error ? err.message : String(err);
    } finally {
      dictionaryBusy = false;
    }
  }
</script>

<div class="grid gap-3 p-4" data-testid="readability-settings">
  {#if error}
    <p class="text-error text-sm" data-testid="readability-error">{error}</p>
  {/if}
  <label class="form-control">
    <span class="label-text">Font size ({readability.fontSize}px)</span>
    <input
      class="range range-primary"
      type="range"
      min="12"
      max="32"
      step="1"
      bind:value={readability.fontSize}
      onchange={changed}
    />
  </label>
  <label class="form-control">
    <span class="label-text">Line height ({readability.lineHeight.toFixed(2)})</span>
    <input
      class="range range-primary"
      type="range"
      min="0.8"
      max="2"
      step="0.05"
      bind:value={readability.lineHeight}
      onchange={changed}
    />
  </label>
  <fieldset class="fieldset">
    <legend class="fieldset-legend">Font family</legend>
    <div class="join grid w-full grid-cols-2">
      {#each fontOptions as option (option.id)}
        <button
          type="button"
          class="btn join-item h-auto min-h-12 flex-col py-2"
          class:btn-primary={readability.fontFamily === option.id}
          style:font-family={readerFontStack(option.id)}
          data-testid="font-family-{option.id}"
          onclick={() => setFontFamily(option.id)}
        >
          {option.label}
        </button>
      {/each}
    </div>
  </fieldset>
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
  <div class="divider my-1">Dictionary</div>
  <button type="button" class="btn btn-xs btn-ghost" data-testid="dictionary-status" onclick={refreshDictionaryStatus}>
    Check dictionary status
  </button>
  <button
    type="button"
    class="btn btn-sm btn-outline"
    data-testid="dictionary-reinstall"
    disabled={dictionaryBusy}
    onclick={reinstallDictionary}
  >
    {dictionaryBusy ? "Reinstalling…" : "Reinstall Webster dictionary"}
  </button>
  {#if dictionaryMsg}
    <p class="text-sm opacity-80">{dictionaryMsg}</p>
  {/if}
</div>
