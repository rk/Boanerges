<script lang="ts">
  import { onMount } from "svelte";
  import { AppChrome, BibleService, events, Service, type ProgressPayload } from "./mygo";
  import {
    addVerseToList,
    applyStudySettings,
    copyVerseReference,
    hydrateStudy,
    openComparisonForVerse,
    openCrossReferencesForVerse,
    openDictionaryForWord,
    setColumnCount,
    study,
    toggleColumn,
  } from "./lib/study.svelte.ts";
  import type { ColumnType } from "./lib/study.svelte.ts";
  import TranslationManagerModal from "./components/sidebar/TranslationManagerModal.svelte";
  import { openStudyPrintPreview } from "./lib/printStudy";
  import PrintPreview from "./pages/PrintPreview.svelte";
  import { bindTranslationInstallProgress, openTranslationManager } from "./lib/translations.svelte.ts";
  import Study from "./pages/Study.svelte";
  import Welcome from "./pages/Welcome.svelte";
  import BookChapterFlyout from "./pages/BookChapterFlyout.svelte";
  import VerseActionsFlyout from "./pages/VerseActionsFlyout.svelte";

  const flyoutParam =
    typeof window !== "undefined" ? new URLSearchParams(window.location.search).get("flyout") : null;
  const viewParam = typeof window !== "undefined" ? new URLSearchParams(window.location.search).get("view") : null;

  let ready = $state(false);
  let studyLoaded = $state(import.meta.env.VITE_E2E === "1");
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
      void openStudyPrintPreview();
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
    const offStudy = events.studySettingsChanged.on(() => {
      if (ready) {
        void Service.getStudy().then(applyStudySettings);
      }
    });
    const offVerseMenu = events.verseMenuAction.on((payload) => {
      const { action, bookId, chapter, verse, word } = payload;
      const ch = Number(chapter);
      const v = Number(verse);
      if (action === "comparison") {
        openComparisonForVerse(String(bookId), ch, v);
        return;
      }
      if (action === "cross-references") {
        openCrossReferencesForVerse(String(bookId), ch, v);
        return;
      }
      if (action === "verse-list") {
        void addVerseToList(String(bookId), ch, v);
        return;
      }
      if (action === "copy-reference") {
        void copyVerseReference(String(bookId), ch, v);
        return;
      }
      if (action === "dictionary" && word) {
        openDictionaryForWord(word);
      }
    });
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
      offStudy();
      offVerseMenu();
      offMenu();
    };
  });

  $effect(() => {
    if (!ready || studyLoaded) {
      return;
    }
    void Service.getStudy().then((settings) => {
      hydrateStudy(settings);
      studyLoaded = true;
    });
  });
</script>

{#if flyoutParam === "verse-actions"}
  <VerseActionsFlyout />
{:else if flyoutParam === "book-chapter"}
  <BookChapterFlyout />
{:else if viewParam === "print-preview"}
  <PrintPreview />
{:else}
  <TranslationManagerModal />

  {#if ready && studyLoaded}
    <Study />
  {:else if ready}
    <p class="flex h-screen items-center justify-center opacity-60">Loading study…</p>
  {:else}
    <Welcome {progress} />
  {/if}
{/if}
