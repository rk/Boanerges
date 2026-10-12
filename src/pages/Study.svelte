<script lang="ts">
  import { onMount } from "svelte";
  import { Service } from "../mygo";
  import ColumnLayout from "../components/layout/ColumnLayout.svelte";
  import ReadabilitySettings from "../components/settings/ReadabilitySettings.svelte";
  import StudyToolbar from "../components/toolbar/StudyToolbar.svelte";
  import { hydrateReadability } from "../lib/readability.svelte.ts";
  import { hydrateStudy, study } from "../lib/study.svelte.ts";

  onMount(() => {
    void Service.getStudy().then(hydrateStudy);
    void Service.getReadability().then(hydrateReadability);
  });
</script>

<div class="flex h-screen flex-col bg-base-100 text-base-content" data-testid="study-screen">
  <StudyToolbar />
  <div class="relative min-h-0 flex-1">
    <ColumnLayout />
    {#if study.settingsOpen}
      <aside class="absolute right-0 top-0 z-10 h-full w-80 border-l border-base-300 bg-base-100 shadow-lg">
        <div class="flex items-center justify-between border-b border-base-300 px-4 py-2">
          <h2 class="font-semibold">Readability</h2>
          <button type="button" class="btn btn-xs btn-ghost" onclick={() => (study.settingsOpen = false)}>Close</button>
        </div>
        <ReadabilitySettings />
      </aside>
    {/if}
  </div>
</div>
