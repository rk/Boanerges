<script lang="ts">
  import { onMount } from "svelte";
  import { BibleService, events, type ProgressPayload } from "./mygo";
  import Study from "./pages/Study.svelte";
  import Welcome from "./pages/Welcome.svelte";

  let ready = $state(false);
  let progress = $state({ step: "starting", percent: 0 });

  function applyStatus(p: ProgressPayload) {
    progress = { step: p.step ?? progress.step, percent: p.percent ?? progress.percent };
    if (p.step === "ready" || p.install_status === "ready") {
      ready = true;
    }
  }

  onMount(() => {
    const off = events.translationInstallProgress.on(applyStatus);
    void BibleService.asvInstallStatus().then(applyStatus);
    return off;
  });
</script>

{#if ready}
  <Study />
{:else}
  <Welcome {progress} />
{/if}
