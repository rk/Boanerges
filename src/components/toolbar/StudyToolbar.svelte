<script lang="ts">
  import { BibleService, TranslationService } from "../../mygo";
  import { setLocation, study } from "../../lib/study.svelte.ts";

  let books = $state<{ id: string; name: string; chapters: number }[]>([]);
  let translations = $state<{ id: string; abbrev: string }[]>([]);

  $effect(() => {
    void BibleService.books(study.translationId).then((resp) => {
      books = resp.books.map((b) => ({ id: b.id, name: b.name, chapters: b.chapters }));
    });
  });

  $effect(() => {
    void TranslationService.installed().then((resp) => {
      translations = resp.translations.map((t) => ({ id: t.id, abbrev: t.abbrev }));
    });
  });

  const maxChapter = $derived(books.find((b) => b.id === study.bookId)?.chapters ?? 1);

  function onBookChange(e: Event) {
    setLocation((e.currentTarget as HTMLSelectElement).value, 1);
  }

  function onChapterChange(e: Event) {
    setLocation(study.bookId, Number((e.currentTarget as HTMLSelectElement).value));
  }

  function onTranslationChange(e: Event) {
    study.translationId = (e.currentTarget as HTMLSelectElement).value;
  }

  function prevChapter() {
    if (study.chapter > 1) setLocation(study.bookId, study.chapter - 1);
  }

  function nextChapter() {
    if (study.chapter < maxChapter) setLocation(study.bookId, study.chapter + 1);
  }
</script>

<header class="navbar bg-base-200 border-b border-base-300 px-3 gap-2">
  <div class="flex flex-wrap items-center gap-2">
    <select class="select select-sm select-bordered" data-testid="translation-select" value={study.translationId} onchange={onTranslationChange}>
      {#each translations as t}
        <option value={t.id}>{t.abbrev}</option>
      {/each}
    </select>
    <select class="select select-sm select-bordered" data-testid="book-select" value={study.bookId} onchange={onBookChange}>
      {#each books as book}
        <option value={book.id}>{book.name}</option>
      {/each}
    </select>
    <div class="join">
      <button type="button" class="btn btn-sm join-item" data-testid="chapter-prev" onclick={prevChapter} disabled={study.chapter <= 1}>‹</button>
      <select class="select select-sm select-bordered join-item w-20" data-testid="chapter-select" value={String(study.chapter)} onchange={onChapterChange}>
        {#each Array.from({ length: maxChapter }, (_, i) => i + 1) as ch}
          <option value={String(ch)}>{ch}</option>
        {/each}
      </select>
      <button type="button" class="btn btn-sm join-item" data-testid="chapter-next" onclick={nextChapter} disabled={study.chapter >= maxChapter}>›</button>
    </div>
  </div>
  <div class="flex-1"></div>
  <button type="button" class="btn btn-sm btn-ghost" onclick={() => (study.settingsOpen = !study.settingsOpen)}>Settings</button>
</header>
