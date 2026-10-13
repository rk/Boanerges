<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { DictionaryService } from "../../mygo";
  import { parseDictionaryLookup, type DictionaryLookupResult } from "../../lib/dictionaryResponse";
  import { normalizeDictionaryWord } from "../../lib/normalizeDictionaryWord";
  import { study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let query = $state("");
  let suggestions = $state<string[]>([]);
  let result = $state<DictionaryLookupResult | null>(null);
  let loading = $state(false);
  let error = $state("");

  async function suggest() {
    const resp = (await DictionaryService.suggest(query, 10)) as { suggestions?: string[]; status?: string };
    if (resp.status === "importing") {
      suggestions = [];
      return;
    }
    suggestions = resp.suggestions ?? [];
  }

  $effect(() => {
    const pending = study.dictionaryQuery;
    if (!pending) {
      return;
    }
    study.dictionaryQuery = null;
    void lookup(pending);
  });

  async function lookup(raw: string) {
    const word = normalizeDictionaryWord(raw);
    if (word.length < 2) {
      result = null;
      error = "";
      return;
    }
    query = word;
    loading = true;
    error = "";
    try {
      result = parseDictionaryLookup(await DictionaryService.lookup(word), word);
    } catch (err: unknown) {
      result = null;
      error = err instanceof Error ? err.message : "Lookup failed";
    } finally {
      loading = false;
    }
  }
</script>

<div class="flex h-full flex-col">
  <ColumnHeader contentType="dictionary" {slotIndex} showViewSelector>
    <input
      class="input input-bordered input-sm w-full"
      data-testid="dictionary-input"
      bind:value={query}
      oninput={suggest}
      onkeydown={(e) => e.key === "Enter" && lookup(query)}
    />
  </ColumnHeader>
  {#if suggestions.length}
    <div class="flex flex-wrap gap-1 p-2">
      {#each suggestions as word}
        <button type="button" class="btn btn-xs btn-outline" onclick={() => lookup(word)}>{word}</button>
      {/each}
    </div>
  {/if}
  <article class="flex-1 overflow-y-auto p-4 text-sm">
    {#if loading}
      <p class="opacity-60">Looking up…</p>
    {:else if error}
      <p class="text-error">{error}</p>
    {:else if result?.found}
      <h3 class="font-semibold">{result.word}</h3>
      {#each result.variants ?? [] as variant}
        <p class="mt-2 opacity-70">{variant.partOfSpeechExpanded}</p>
        <ul class="list-disc pl-5">
          {#each variant.definitions as def}
            <li>{def}</li>
          {/each}
        </ul>
      {/each}
    {:else if result?.message}
      <p class="opacity-70">{result.message}</p>
    {/if}
  </article>
</div>
