<script lang="ts">
  import type { RunState } from "../agentic/client";

  let {
    run,
    startLabel = "Start run",
    onStart,
    onPause,
    onResume,
    startDisabled = false,
    startTestId,
    showStart = true,
  }: {
    run: RunState | undefined;
    startLabel?: string;
    onStart: () => void | Promise<void>;
    onPause: () => void | Promise<void>;
    onResume: () => void | Promise<void>;
    startDisabled?: boolean;
    startTestId?: string;
    showStart?: boolean;
  } = $props();

  const status = $derived(run?.status ?? "idle");
  const badgeClass = $derived(
    status === "running"
      ? "badge-primary"
      : status === "paused"
        ? "badge-warning"
        : status === "failed"
          ? "badge-error"
          : "badge-ghost",
  );
</script>

<div class="flex flex-wrap items-center gap-2">
  {#if showStart}
    <button
      type="button"
      class="btn btn-sm btn-primary"
      disabled={startDisabled || status === "running" || !run}
      data-testid={startTestId}
      onclick={() => void onStart()}
    >
      {startLabel}
    </button>
  {/if}
  <button type="button" class="btn btn-sm btn-outline" disabled={status !== "running"} onclick={() => void onPause()}>
    Pause
  </button>
  <button type="button" class="btn btn-sm btn-outline" disabled={status !== "paused"} onclick={() => void onResume()}>
    Resume
  </button>
  <span class="badge badge-sm {badgeClass}">{status}</span>
</div>
