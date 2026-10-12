<script lang="ts">
  import { TranslationService } from "../../mygo";
  import { bible, fetchBooksForTranslations } from "../../lib/bible.svelte.ts";
  import { schedulePersistStudy, setLocation, study } from "../../lib/study.svelte.ts";
  import BookChapterPicker from "./BookChapterPicker.svelte";

  let translations = $state<{ id: string; abbrev: string }[]>([]);

  $effect(() => {
    void fetchBooksForTranslations([study.translationId]);
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

  <div class="flex-1"></div>
  <button type="button" class="btn btn-sm btn-ghost" onclick={() => (study.settingsOpen = !study.settingsOpen)}>Settings</button>
</header>
