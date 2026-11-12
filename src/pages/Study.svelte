<script lang="ts">
  import { onMount } from "svelte";
  import { Service } from "../mygo";
  import ColumnLayout from "../components/layout/ColumnLayout.svelte";
  import ReadabilitySettings from "../components/settings/ReadabilitySettings.svelte";
  import StudyToolbar from "../components/toolbar/StudyToolbar.svelte";
  import { refreshInstalledTranslations } from "../lib/bible.svelte.ts";
  import { hydrateReadability } from "../lib/readability.svelte.ts";
  import { hydrateStudy, study } from "../lib/study.svelte.ts";

  onMount(() => {
    void refreshInstalledTranslations();
    void Service.getStudy().then(hydrateStudy);
    void Service.getReadability().then(hydrateReadability);
  });
</script>

<div class="flex h-screen flex-col bg-base-200 text-base-content study-workspace" data-testid="study-screen">
  <StudyToolbar />
  <div class="relative min-h-0 flex-1">
    <ColumnLayout />
  </div>
</div>

<div class="drawer drawer-end z-20">
  <input id="study-settings-drawer" type="checkbox" class="drawer-toggle" bind:checked={study.settingsOpen} />
  <div class="drawer-content"></div>
  <div class="drawer-side">
    <label
      for="study-settings-drawer"
      class="drawer-overlay"
      aria-label="Close settings"
      onclick={() => (study.settingsOpen = false)}
    ></label>
    <aside class="bg-base-100 flex h-full w-80 max-w-[90vw] flex-col border-l border-base-300">
      <div class="flex items-center justify-between border-b border-base-300 px-4 py-3">
        <h2 class="text-sm font-semibold">Readability</h2>
        <label for="study-settings-drawer" class="btn btn-sm btn-ghost">Close</label>
      </div>
      <div class="min-h-0 flex-1 overflow-y-auto">
        <ReadabilitySettings />
      </div>
    </aside>
  </div>
</div>
