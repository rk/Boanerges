<script lang="ts">
  import { onMount } from "svelte";
  import { AppChrome, BibleService, events, type ProgressPayload } from "./mygo";
  import { setColumnCount, study, toggleColumn } from "./lib/study.svelte.ts";
  import type { ColumnType } from "./lib/study.svelte.ts";
  import TranslationManagerModal from "./components/sidebar/TranslationManagerModal.svelte";
  import PrintOptionsModal from "./components/print/PrintOptionsModal.svelte";
  import { bindTranslationInstallProgress, openTranslationManager } from "./lib/translations.svelte.ts";
  import Study from "./pages/Study.svelte";
  import Welcome from "./pages/Welcome.svelte";

  let printOpen = $state(false);

  let ready = $state(false);
  let progress = $state({ step: "starting", percent: 0, error: "" });

  function applyStatus(p: ProgressPayload) {
    progress = {
      step: p.step ?? progress.step,
      percent: p.percent ?? progress.percent,
      error: p.error ?? "",
    };
    if (p.step === "ready" || p.install_status === "ready") {
      ready = true;
    }
  }

  function syncMenu() {
    if (!ready || import.meta.env.VITE_E2E === "1") {
      return;
    }
    void AppChrome.syncMenuState(study.columnCount, study.scrollSync);
  }

  function handleMenu(id: string) {
    if (id === "file.manage-translations") {
      void openTranslationManager();
      return;
    }
    if (id === "study.feedback") {
      void AppChrome.openFeedback();
      return;
    }
    if (id === "view.devtools") {
      AppChrome.openDevTools();
      return;
    }
    if (id === "study.print") {
      printOpen = true;
      return;
    }
    if (id === "view.settings") {
      study.settingsOpen = !study.settingsOpen;
      syncMenu();
      return;
    }
    if (id === "view.scroll-sync") {
      study.scrollSync = !study.scrollSync;
      syncMenu();
      return;
    }
    if (id.startsWith("view.columns.")) {
      setColumnCount(Number(id.slice(-1)) as 1 | 2 | 3);
      syncMenu();
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

  $effect(() => {
    study.columnCount;
    study.scrollSync;
    syncMenu();
  });

  onMount(() => {
    if (import.meta.env.VITE_E2E === "1") {
      window.__openTranslationManager = () => openTranslationManager();
    }
    const offProgress = events.translationInstallProgress.on(applyStatus);
    const offTranslationProgress = bindTranslationInstallProgress();
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
      offTranslationProgress();
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
