<script lang="ts">
  import { onMount } from "svelte";
  import { BibleService, events, type ProgressPayload } from "./mygo";
  import { setColumnCount, study, toggleColumn } from "./lib/study.svelte.ts";
  import type { ColumnType } from "./lib/study.svelte.ts";
  import TranslationManagerModal from "./components/sidebar/TranslationManagerModal.svelte";
  import PrintOptionsModal from "./components/print/PrintOptionsModal.svelte";
  import { openTranslationManager } from "./lib/translations.svelte.ts";
  import Study from "./pages/Study.svelte";
  import Welcome from "./pages/Welcome.svelte";

  let printOpen = $state(false);

  let ready = $state(false);
  let progress = $state({ step: "starting", percent: 0 });

  function applyStatus(p: ProgressPayload) {
    progress = { step: p.step ?? progress.step, percent: p.percent ?? progress.percent };
    if (p.step === "ready" || p.install_status === "ready") {
      ready = true;
    }
  }

  function handleMenu(id: string) {
    if (id === "file.manage-translations") {
      void openTranslationManager();
      return;
    }
    if (id === "study.print") {
      printOpen = true;
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
      "study.notes": "notes",
      "study.scribe": "scribe",
      "study.search": "search",
      "study.cross-references": "cross-references",
      "study.dictionary": "dictionary",
      "study.comparison": "comparison",
      "study.verse-list": "verse-list",
    };
    if (map[id]) toggleColumn(map[id]);
  }

  onMount(() => {
    if (import.meta.env.VITE_E2E === "1") {
      window.__openTranslationManager = () => openTranslationManager();
    }
    const offProgress = events.translationInstallProgress.on(applyStatus);
    const offMenu = events.studyMenu.on(handleMenu);
    const welcomeQuery =
      typeof window !== "undefined" && new URLSearchParams(window.location.search).get("e2eWelcome") === "1";
    const mode = import.meta.env.VITE_E2E;
    if (welcomeQuery) {
      // welcome-flow e2e drives progress via mock events
    } else if (mode === "1") {
      ready = true;
    } else {
      void BibleService.asvInstallStatus().then(applyStatus);
    }
    return () => {
      offProgress();
      offMenu();
    };
  });
</script>

<TranslationManagerModal />
<PrintOptionsModal bind:open={printOpen} />

{#if ready}
  <Study />
{:else}
  <Welcome {progress} />
{/if}
