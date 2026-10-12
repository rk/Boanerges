<script lang="ts">
  import { onMount } from "svelte";
  import { AppChrome, BibleService, events, type ProgressPayload } from "./mygo";
  import { setColumnCount, study, toggleColumn } from "./lib/study.svelte.ts";
  import type { ColumnType } from "./lib/study.svelte.ts";
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

  function handleMenu(id: string) {
    if (id === "study.print") {
      void AppChrome.print();
      return;
    }
    if (id === "view.settings") {
      study.settingsOpen = !study.settingsOpen;
      return;
    }
    if (id === "view.scroll-sync") {
      study.scrollSync = !study.scrollSync;
      return;
    }
    if (id.startsWith("view.columns.")) {
      setColumnCount(Number(id.slice(-1)) as 1 | 2 | 3);
      return;
    }
    const map: Record<string, ColumnType> = {
      "study.search": "search",
      "study.cross-references": "cross-references",
      "study.dictionary": "dictionary",
      "study.comparison": "comparison",
      "study.verse-list": "verse-list",
    };
    if (map[id]) toggleColumn(map[id]);
  }

  onMount(() => {
    const offProgress = events.translationInstallProgress.on(applyStatus);
    const offMenu = events.studyMenu.on(handleMenu);
    void BibleService.asvInstallStatus().then(applyStatus);
    return () => {
      offProgress();
      offMenu();
    };
  });
</script>

{#if ready}
  <Study />
{:else}
  <Welcome {progress} />
{/if}
