<script lang="ts">
  import { DictionaryService } from "../../mygo";

  let query = $state("");
  let suggestions = $state<string[]>([]);
  let result = $state<{ found: boolean; word: string; variants?: { definitions: string[]; partOfSpeechExpanded: string }[]; message?: string } | null>(null);

  async function suggest() {
    const resp = (await DictionaryService.suggest(query, 10)) as { suggestions?: string[] };
    suggestions = resp.suggestions ?? [];
  }

  async function lookup(word: string) {
    query = word;
    result = (await DictionaryService.lookup(word)) as typeof result;
  }
</script>

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 p-3">
    <input
      class="input input-bordered input-sm w-full"
      data-testid="dictionary-input"
      bind:value={query}
      oninput={suggest}
      onkeydown={(e) => e.key === "Enter" && lookup(query)}
    />
  </div>
  {#if suggestions.length}
    <div class="flex flex-wrap gap-1 p-2">
      {#each suggestions as word}
        <button type="button" class="btn btn-xs btn-outline" onclick={() => lookup(word)}>{word}</button>
      {/each}
    </div>
  {/if}
  <article class="flex-1 overflow-y-auto p-4 text-sm">
    {#if result?.found}
      <h3 class="font-semibold">{result.word}</h3>
      {#each result.variants ?? [] as variant}
        <p class="mt-2 opacity-70">{variant.partOfSpeechExpanded}</p>
        <ul class="list-disc pl-5">
          {#each variant.definitions as def}
            <li>{def}</li>
          {/each}
        </ul>
      {/each}
    {:else if result}
      <p>{result.message}</p>
    {/if}
  </article>
</div>
