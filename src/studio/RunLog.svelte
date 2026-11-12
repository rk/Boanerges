<script lang="ts">
  import { agentic, type RunState } from "../agentic/client";

  let {
    run,
    folder = "",
  }: {
    run: RunState;
    folder?: string;
  } = $props();

  let showRaw = $state(false);
  let raw = $state("");
  let rawLoading = $state(false);

  function formatWhen(iso: string): string {
    const d = new Date(iso);
    if (Number.isNaN(d.getTime())) {
      return iso;
    }
    return d.toLocaleTimeString(undefined, { hour: "2-digit", minute: "2-digit", second: "2-digit" });
  }

  async function loadRaw(): Promise<void> {
    if (!folder) {
      return;
    }
    rawLoading = true;
    try {
      raw = await agentic.readRawRunLog(folder, 300_000);
    } finally {
      rawLoading = false;
    }
  }

  async function toggleRaw(): Promise<void> {
    showRaw = !showRaw;
    if (showRaw) {
      await loadRaw();
    }
  }
</script>

<div
  class="border-base-300 bg-base-100 flex min-h-48 flex-1 flex-col rounded-lg border shadow-sm"
  data-testid="studio-run-log"
>
  <div class="border-base-300/80 flex flex-wrap items-center justify-between gap-2 border-b px-4 py-2">
    <span class="text-sm font-semibold">Run log</span>
    <div class="flex items-center gap-2">
      {#if folder}
        <button type="button" class="btn btn-ghost btn-xs" data-testid="studio-raw-log-toggle" onclick={() => void toggleRaw()}>
          {showRaw ? "Summary" : "Raw debug"}
        </button>
        {#if showRaw}
          <button type="button" class="btn btn-ghost btn-xs" disabled={rawLoading} onclick={() => void loadRaw()}>Reload raw</button>
        {/if}
      {/if}
      <span
        class="badge badge-sm {run.status === 'running'
          ? 'badge-primary'
          : run.status === 'failed'
            ? 'badge-error'
            : 'badge-ghost'}"
      >
        {run.status}
      </span>
    </div>
  </div>
  {#if run.error}
    <p class="text-error border-base-300/80 border-b px-4 py-2 text-xs">{run.error}</p>
  {/if}
  {#if showRaw}
    <pre
      class="bg-base-200/40 text-base-content/90 max-h-96 min-h-48 flex-1 overflow-auto p-3 font-mono text-[0.65rem] leading-snug whitespace-pre-wrap"
      data-testid="studio-raw-log"
    >{#if rawLoading}Loading…{:else if raw}{raw}{:else}No run.raw.log yet — start a run.{/if}</pre>
    <p class="text-base-content/50 border-base-300/80 border-t px-4 py-1 text-[0.65rem]">Also on disk: <code>run.raw.log</code></p>
  {:else}
    <ol class="min-h-0 flex-1 list-none space-y-0 overflow-y-auto px-4 py-2 text-xs leading-relaxed">
      {#each run.steps as step, i (step.at + String(i))}
        <li class="border-base-300/50 border-b py-2 last:border-b-0">
          <div class="text-base-content/50 mb-0.5 font-mono text-[0.65rem]">
            {formatWhen(step.at)}
            <span class="text-base-content/80 mx-1.5 font-sans font-medium">{step.role}</span>
          </div>
          <p class="text-base-content/90 whitespace-pre-wrap break-words">{step.summary}</p>
        </li>
      {:else}
        <li class="text-base-content/50 py-6 text-center">No steps yet — start a run to see tool calls here.</li>
      {/each}
    </ol>
  {/if}
</div>
