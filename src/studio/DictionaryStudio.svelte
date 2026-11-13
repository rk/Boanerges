<script lang="ts">
  import { onMount } from "svelte";
  import { agentic, agenticEvents, type EntrySummary, type ProjectSummary } from "../agentic/client";
  import OllamaBar from "./OllamaBar.svelte";
  import RunControls from "./RunControls.svelte";
  import RunLog from "./RunLog.svelte";

  let {
    folder = $bindable(""),
    summary = $bindable<ProjectSummary | null>(null),
  }: {
    folder?: string;
    summary?: ProjectSummary | null;
  } = $props();

  let sources = $state<string[]>([]);
  let entries = $state<EntrySummary[]>([]);
  let instruction = $state(
    "Work through every source file to completion. read_source with args.name (basename only), limit 8000. Extract real dictionary headwords from the chunk text only — never from these instructions. At least 4 upsert_entry per chunk before the next read.",
  );
  let model = $state("llama3.2");
  let ollamaReady = $state(false);

  async function refreshLists(): Promise<void> {
    if (!folder) {
      return;
    }
    sources = await agentic.listSources(folder);
    entries = await agentic.listEntries(folder);
  }

  async function refresh(): Promise<void> {
    if (!folder) {
      return;
    }
    try {
      summary = await agentic.loadProjectSummary(folder);
      if (summary.manifest.kind !== "dictionary") {
        summary = null;
        folder = "";
        return;
      }
      await refreshLists();
      model = summary.manifest.model;
    } catch {
      summary = null;
      sources = [];
      entries = [];
      folder = "";
    }
  }

  async function pickFolder(): Promise<void> {
    const path = await agentic.pickProjectFolder();
    if (!path) {
      return;
    }
    await agentic.initProject(path, "dictionary");
    folder = path;
    await refresh();
    if (model) {
      await agentic.setProjectModel(folder, model);
      await refresh();
    }
  }

  async function saveModel(name: string): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.setProjectModel(folder, name);
    await refresh();
  }

  async function analyzeFormats(): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.analyzeSourceFormats(folder);
    await refresh();
  }

  async function start(): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.startDictionaryRun(folder, instruction);
    await refresh();
  }

  async function pause(): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.pauseRun(folder);
    await refresh();
  }

  async function resume(): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.resumeRun(folder);
    await refresh();
  }

  onMount(() => {
    const off = agenticEvents.runStateChanged.on(() => {
      void refresh();
    });
    return () => off();
  });

  $effect(() => {
    folder;
    void refresh();
  });

  $effect(() => {
    if (!folder || summary?.run.status !== "running") {
      return;
    }
    const id = setInterval(() => {
      void refreshLists();
      void agentic.loadProjectSummary(folder).then((s) => {
        summary = s;
      });
    }, 2000);
    return () => clearInterval(id);
  });
</script>

<div class="mx-auto flex min-h-0 w-full max-w-5xl flex-1 flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <button type="button" class="btn btn-sm" data-testid="dict-pick-folder" onclick={pickFolder}>
      {folder ? "Change folder" : "Choose project folder"}
    </button>
    {#if folder}
      <p class="text-base-content/60 min-w-0 flex-1 truncate font-mono text-xs" title={folder}>{folder}</p>
    {/if}
  </div>

  <OllamaBar bind:model bind:ready={ollamaReady} disabled={!folder} onModelPick={saveModel} />

  {#if !folder}
    <div class="card bg-base-200/50 border-base-300 border">
      <div class="card-body items-center text-center">
        <h2 class="card-title text-base">Dictionary project</h2>
        <p class="text-base-content/70 max-w-md text-sm">
          Pick a folder on disk. Drop public-domain sources into <code class="text-xs">sources/</code>, then run the
          distiller with a local Ollama model.
        </p>
        <button type="button" class="btn btn-primary btn-sm mt-2" onclick={pickFolder}>Choose project folder</button>
      </div>
    </div>
  {:else}
    <div class="grid min-h-0 flex-1 gap-4 lg:grid-cols-[minmax(0,18rem)_minmax(0,1fr)]">
      <aside class="card bg-base-100 border-base-300 flex flex-col gap-4 border shadow-sm">
        <div class="card-body gap-4 p-4">
          <label class="form-control gap-1">
            <span class="text-xs font-medium">Standing instruction</span>
            <textarea
              class="textarea textarea-bordered textarea-sm min-h-28 w-full text-sm leading-snug"
              bind:value={instruction}
            ></textarea>
          </label>
          <button
            type="button"
            class="btn btn-sm btn-outline w-full"
            disabled={!ollamaReady || !model || summary?.run.status === "running"}
            data-testid="dict-analyze-formats"
            onclick={() => void analyzeFormats()}
          >
            Analyze source formats
          </button>
          <p class="text-base-content/50 text-[0.65rem] leading-snug">
            Reads first 10k per source (parallel subagents) into <code>source-guides.json</code>. Runs automatically on Start too.
          </p>
          <RunControls
            run={summary?.run}
            onStart={start}
            onPause={pause}
            onResume={resume}
            startDisabled={!ollamaReady || !model}
            startTestId="dict-start"
          />
        </div>
      </aside>

      <div class="flex min-h-0 flex-1 flex-col gap-4">
        <div class="grid gap-4 sm:grid-cols-2">
          <section class="card bg-base-100 border-base-300 border shadow-sm">
            <div class="card-body p-4">
              <div class="mb-2 flex items-start justify-between gap-2">
                <div>
                  <h3 class="text-sm font-semibold">Sources <span class="text-base-content/50 font-normal">({sources.length})</span></h3>
                  <p class="text-base-content/60 mt-0.5 text-xs">
                    Drop UTF-8 <code class="text-[0.7rem]">.txt</code> or <code class="text-[0.7rem]">.md</code> in
                    <code class="text-[0.7rem]">sources/</code> — not PDF.
                  </p>
                </div>
                <button
                  type="button"
                  class="btn btn-ghost btn-square btn-sm shrink-0"
                  title="Refresh source list"
                  aria-label="Refresh source list"
                  data-testid="dict-refresh-sources"
                  onclick={() => void refresh()}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h5M20 20v-5h-5M20 9a8 8 0 0 0-14.5-4M4 15a8 8 0 0 0 14.5 4" />
                  </svg>
                </button>
              </div>
              <ul class="bg-base-200/40 max-h-36 overflow-y-auto rounded-md border border-base-300/80 p-2 text-xs leading-relaxed">
                {#each sources as name (name)}
                  <li class="truncate py-0.5" title={name}>{name}</li>
                {:else}
                  <li class="text-base-content/50">No source files yet</li>
                {/each}
              </ul>
            </div>
          </section>
          <section class="card bg-base-100 border-base-300 border shadow-sm">
            <div class="card-body p-4">
              <div class="mb-2 flex items-center justify-between gap-2">
                <h3 class="text-sm font-semibold">Entries <span class="text-base-content/50 font-normal">({entries.length})</span></h3>
                <button
                  type="button"
                  class="btn btn-ghost btn-square btn-sm shrink-0"
                  title="Refresh entries list"
                  aria-label="Refresh entries list"
                  data-testid="dict-refresh-entries"
                  onclick={() => void refreshLists()}
                >
                  <svg xmlns="http://www.w3.org/2000/svg" class="size-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h5M20 20v-5h-5M20 9a8 8 0 0 0-14.5-4M4 15a8 8 0 0 0 14.5 4" />
                  </svg>
                </button>
              </div>
              <ul class="bg-base-200/40 max-h-36 overflow-y-auto rounded-md border border-base-300/80 p-2 text-xs leading-relaxed">
                {#each entries as e (e.file)}
                  <li class="truncate py-0.5">{e.headword}</li>
                {:else}
                  <li class="text-base-content/50">No entries yet</li>
                {/each}
              </ul>
            </div>
          </section>
        </div>
        {#if summary}
          <RunLog run={summary.run} {folder} />
        {/if}
      </div>
    </div>
  {/if}
</div>
