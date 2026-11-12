<script lang="ts">
  import { onMount } from "svelte";
  import { agentic } from "../agentic/client";

  let {
    model = $bindable(""),
    ready = $bindable(false),
    disabled = false,
    onModelPick,
  }: {
    model?: string;
    ready?: boolean;
    disabled?: boolean;
    onModelPick?: (name: string) => void | Promise<void>;
  } = $props();

  let models = $state<string[]>([]);
  let loading = $state(false);
  let error = $state("");

  const modelOptions = $derived.by(() => {
    const names = [...models];
    if (model && !names.includes(model)) {
      names.unshift(model);
    }
    return names;
  });

  async function refresh(): Promise<void> {
    loading = true;
    error = "";
    try {
      const h = await agentic.pingOllama();
      ready = h.ok;
      models = h.models ?? [];
      if (!h.ok) {
        error = h.error ?? "Ollama unreachable";
        return;
      }
      if (!model && models.length > 0) {
        model = models[0]!;
      }
    } finally {
      loading = false;
    }
  }

  async function pick(name: string): Promise<void> {
    model = name;
    await onModelPick?.(name);
  }

  onMount(() => {
    void refresh();
  });
</script>

<div class="flex flex-wrap items-end gap-3">
  <div class="flex min-w-0 flex-1 flex-col gap-1 sm:max-w-md">
    <span class="text-xs font-medium">Ollama model</span>
    <div class="flex gap-2">
      <select
        class="select select-bordered select-sm min-w-0 flex-1"
        {disabled}
        value={model}
        onchange={(e) => void pick(e.currentTarget.value)}
        aria-label="Ollama model"
      >
        {#if modelOptions.length === 0}
          <option value="">No models — refresh after starting Ollama</option>
        {:else}
          {#each modelOptions as name (name)}
            <option value={name}>{name}</option>
          {/each}
        {/if}
      </select>
      <button
        type="button"
        class="btn btn-sm btn-ghost shrink-0"
        disabled={loading}
        title="Reload model list from Ollama"
        onclick={() => void refresh()}
      >
        {loading ? "…" : "Refresh"}
      </button>
    </div>
  </div>
  <div class="flex items-center gap-2 pb-1 text-xs">
    {#if loading}
      <span class="loading loading-spinner loading-xs"></span>
      <span class="text-base-content/60">Checking Ollama…</span>
    {:else if ready}
      <span class="status status-success status-xs" aria-hidden="true"></span>
      <span class="text-base-content/70">{models.length} model{models.length === 1 ? "" : "s"} installed</span>
    {:else}
      <span class="status status-error status-xs" aria-hidden="true"></span>
      <span class="text-error">{error || "Ollama not ready"}</span>
    {/if}
  </div>
</div>
