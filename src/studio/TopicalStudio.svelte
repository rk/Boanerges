<script lang="ts">
  import { onMount } from "svelte";
  import { agentic, agenticEvents, type ProjectSummary, type TopicSummary } from "../agentic/client";
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

  let topics = $state<TopicSummary[]>([]);
  let instruction = $state("Identify topics, reuse existing topics when titles match, group refs with short labels.");
  let model = $state("llama3.2");
  let ollamaReady = $state(false);
  let refBook = $state("gen");
  let refChapter = $state(1);
  let refVerse = $state(1);
  let topicName = $state("Faith");

  async function refresh(): Promise<void> {
    if (!folder) {
      return;
    }
    try {
      summary = await agentic.loadProjectSummary(folder);
      if (summary.manifest.kind !== "topical") {
        summary = null;
        folder = "";
        return;
      }
      topics = await agentic.listTopics(folder);
      model = summary.manifest.model;
    } catch {
      summary = null;
      topics = [];
      folder = "";
    }
  }

  async function pickFolder(): Promise<void> {
    const path = await agentic.pickProjectFolder();
    if (!path) {
      return;
    }
    await agentic.initProject(path, "topical");
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

  async function startFromRef(): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.startTopicalFromReference(folder, refBook, refChapter, refVerse, instruction);
    await refresh();
  }

  async function startFromTopic(): Promise<void> {
    if (!folder) {
      return;
    }
    await agentic.startTopicalFromTopic(folder, topicName, instruction);
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
</script>

<div class="mx-auto flex min-h-0 w-full max-w-5xl flex-1 flex-col gap-4">
  <div class="flex flex-wrap items-center justify-between gap-3">
    <button type="button" class="btn btn-sm" data-testid="topical-pick-folder" onclick={pickFolder}>
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
        <h2 class="card-title text-base">Topical project</h2>
        <p class="text-base-content/70 max-w-md text-sm">
          Build grouped scripture references by verse or by topic name. Results land in <code class="text-xs">topics/</code
          > on disk.
        </p>
        <button type="button" class="btn btn-primary btn-sm mt-2" onclick={pickFolder}>Choose project folder</button>
      </div>
    </div>
  {:else}
    <div class="grid min-h-0 flex-1 gap-4 lg:grid-cols-[minmax(0,18rem)_minmax(0,1fr)]">
      <aside class="card bg-base-100 border-base-300 border shadow-sm">
        <div class="card-body gap-4 p-4">
          <label class="form-control gap-1">
            <span class="text-xs font-medium">Standing instruction</span>
            <textarea
              class="textarea textarea-bordered textarea-sm min-h-24 w-full text-sm leading-snug"
              bind:value={instruction}
            ></textarea>
          </label>

          <div class="border-base-300/80 flex flex-col gap-3 border-t pt-3">
            <p class="text-xs font-medium">Start job</p>
            <fieldset class="fieldset p-0">
              <legend class="fieldset-legend text-xs">From reference</legend>
              <div class="grid grid-cols-[minmax(0,1fr)_4rem_4rem] gap-2">
                <input class="input input-bordered input-sm" bind:value={refBook} placeholder="book id" aria-label="Book id" />
                <input class="input input-bordered input-sm" type="number" bind:value={refChapter} aria-label="Chapter" min="1" />
                <input class="input input-bordered input-sm" type="number" bind:value={refVerse} aria-label="Verse" min="1" />
              </div>
              <button
                type="button"
                class="btn btn-sm btn-primary mt-2 w-full"
                data-testid="topical-start-ref"
                disabled={!ollamaReady || !model || summary?.run.status === "running"}
                onclick={() => void startFromRef()}
              >
                Run from reference
              </button>
            </fieldset>
            <fieldset class="fieldset p-0">
              <legend class="fieldset-legend text-xs">From topic name</legend>
              <input class="input input-bordered input-sm w-full" bind:value={topicName} aria-label="Topic" />
              <button
                type="button"
                class="btn btn-sm btn-primary mt-2 w-full"
                data-testid="topical-start-name"
                disabled={!ollamaReady || !model || summary?.run.status === "running"}
                onclick={() => void startFromTopic()}
              >
                Run from topic
              </button>
            </fieldset>
          </div>

          <RunControls run={summary?.run} showStart={false} onStart={startFromRef} onPause={pause} onResume={resume} />
        </div>
      </aside>

      <div class="flex min-h-0 flex-1 flex-col gap-4">
        <section class="card bg-base-100 border-base-300 border shadow-sm">
          <div class="card-body p-4">
            <h3 class="text-sm font-semibold">Topics <span class="text-base-content/50 font-normal">({topics.length})</span></h3>
            <ul class="bg-base-200/40 mt-2 max-h-48 overflow-y-auto rounded-md border border-base-300/80 p-2 text-xs leading-relaxed">
              {#each topics as t (t.slug)}
                <li class="truncate py-0.5">{t.title}</li>
              {:else}
                <li class="text-base-content/50">No topics yet</li>
              {/each}
            </ul>
          </div>
        </section>
        {#if summary}
          <RunLog run={summary.run} {folder} />
        {/if}
      </div>
    </div>
  {/if}
</div>
