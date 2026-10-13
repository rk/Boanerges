<script lang="ts">
  import { TranslationService } from "../../mygo";
  import { bible, fetchBooksForTranslations } from "../../lib/bible.svelte.ts";
  import ScrollSyncToggle from "../reader/ScrollSyncToggle.svelte";
  import { schedulePersistStudy, setColumnCount, setLocation, study } from "../../lib/study.svelte.ts";
  import { activeBibleTranslationIds, bibleColumnCount } from "../../lib/studyLayout";
  import BookChapterPicker from "./BookChapterPicker.svelte";

  const showScrollSync = $derived(bibleColumnCount(study) >= 2);

  let translations = $state<{ id: string; abbrev: string }[]>([]);

  $effect(() => {
    void fetchBooksForTranslations(activeBibleTranslationIds(study));
  });

  $effect(() => {
    void TranslationService.installed().then((resp) => {
      translations = resp.translations.map((t) => ({ id: t.id, abbrev: t.abbrev }));
    });
  });

  const maxChapter = $derived(bible.books.find((b) => b.id === study.bookId)?.chapters ?? 1);

  function onTranslationChange(e: Event) {
    study.translationId = (e.currentTarget as HTMLSelectElement).value;
    void fetchBooksForTranslations([study.translationId]);
    schedulePersistStudy();
  }

  function prevChapter() {
    if (study.chapter > 1) setLocation(study.bookId, study.chapter - 1);
  }

  function nextChapter() {
    if (study.chapter < maxChapter) setLocation(study.bookId, study.chapter + 1);
  }
</script>

<header class="navbar bg-base-200 border-b border-base-300 shrink-0 gap-2 px-3">
  <span class="text-lg font-semibold shrink-0 hidden sm:inline">Boanerges</span>
  <div class="divider divider-horizontal mx-0 hidden h-full sm:flex"></div>

  <select
    class="select select-sm select-bordered"
    data-testid="translation-select"
    value={study.translationId}
    onchange={onTranslationChange}
  >
    {#each translations as t}
      <option value={t.id}>{t.abbrev}</option>
    {/each}
  </select>

  <BookChapterPicker />

  <div class="join">
    <button
      type="button"
      class="btn btn-sm join-item"
      data-testid="chapter-prev"
      onclick={prevChapter}
      disabled={study.chapter <= 1}
      aria-label="Previous chapter"
    >
      ‹
    </button>
    <button
      type="button"
      class="btn btn-sm join-item"
      data-testid="chapter-next"
      onclick={nextChapter}
      disabled={study.chapter >= maxChapter}
      aria-label="Next chapter"
    >
      ›
    </button>
  </div>

  <div class="divider divider-horizontal mx-0 hidden h-full md:flex"></div>

  <div class="join hidden md:flex" data-testid="column-count">
    {#each [1, 2, 3] as count (count)}
      <button
        type="button"
        class="btn btn-sm join-item"
        class:btn-active={study.columnCount === count}
        aria-pressed={study.columnCount === count}
        data-testid="column-count-{count}"
        onclick={() => setColumnCount(count as 1 | 2 | 3)}
      >
        {count}
      </button>
    {/each}
  </div>

  <div class="flex-1"></div>

  {#if showScrollSync}
    <ScrollSyncToggle />
  {/if}

  <button type="button" class="btn btn-sm btn-ghost" onclick={() => (study.settingsOpen = !study.settingsOpen)}>Settings</button>
</header>
