<script lang="ts">
  import { onMount } from "svelte";
  import { agentic, type ProjectSummary } from "../agentic/client";
  import DictionaryStudio from "./DictionaryStudio.svelte";
  import TopicalStudio from "./TopicalStudio.svelte";
  import {
    clearStudioFolder,
    loadStudioFolder,
    saveStudioFolder,
    type StudioProjectKind,
  } from "./projectFolders";

  let tab = $state<"dictionary" | "topical">("dictionary");
  let dictFolder = $state("");
  let topicalFolder = $state("");
  let dictSummary = $state<ProjectSummary | null>(null);
  let topicalSummary = $state<ProjectSummary | null>(null);

  async function restoreFolder(kind: StudioProjectKind, path: string): Promise<string> {
    if (!path) {
      return "";
    }
    try {
      const summary = await agentic.loadProjectSummary(path);
      if (summary.manifest.kind !== kind) {
        clearStudioFolder(kind);
        return "";
      }
      return path;
    } catch {
      clearStudioFolder(kind);
      return "";
    }
  }

  onMount(() => {
    void (async () => {
      dictFolder = await restoreFolder("dictionary", loadStudioFolder("dictionary"));
      topicalFolder = await restoreFolder("topical", loadStudioFolder("topical"));
    })();
  });

  $effect(() => {
    saveStudioFolder("dictionary", dictFolder);
  });

  $effect(() => {
    saveStudioFolder("topical", topicalFolder);
  });
</script>

<div class="bg-base-100 text-base-content flex h-screen min-h-0 flex-col" data-testid="agentic-studio">
  <header class="border-base-300 bg-base-100 flex shrink-0 flex-col gap-2 border-b px-4 py-3 sm:flex-row sm:items-baseline sm:gap-3">
    <h1 class="text-lg font-semibold tracking-tight">Agentic Studio</h1>
    <p class="text-base-content/60 text-sm">Local Ollama · public-domain distillation</p>
  </header>

  <div class="border-base-300 border-b px-4 pt-2" role="tablist">
    <div class="tabs tabs-boxed tabs-sm bg-base-200/60 w-fit max-w-full">
      <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={tab === "dictionary"}
        aria-selected={tab === "dictionary"}
        onclick={() => (tab = "dictionary")}
      >
        Dictionary
      </button>
      <button
        type="button"
        role="tab"
        class="tab"
        class:tab-active={tab === "topical"}
        aria-selected={tab === "topical"}
        onclick={() => (tab = "topical")}
      >
        Topical
      </button>
    </div>
  </div>

  <main class="bg-base-200/30 flex min-h-0 flex-1 flex-col overflow-y-auto p-4 sm:p-6">
    {#if tab === "dictionary"}
      <DictionaryStudio bind:folder={dictFolder} bind:summary={dictSummary} />
    {:else}
      <TopicalStudio bind:folder={topicalFolder} bind:summary={topicalSummary} />
    {/if}
  </main>
</div>
