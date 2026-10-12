<script lang="ts">
  import { BibleService } from "../../mygo";
  import { setLocation, study } from "../../lib/study.svelte.ts";

  let books = $state<{ id: string; name: string; chapters: number }[]>([]);

  $effect(() => {
    void BibleService.books(study.translationId).then((resp) => {
      books = resp.books.map((b) => ({ id: b.id, name: b.name, chapters: b.chapters }));
    });
  });

  function onBookChange(e: Event) {
    const bookId = (e.currentTarget as HTMLSelectElement).value;
    setLocation(bookId, 1);
  }

  function onChapterChange(e: Event) {
    setLocation(study.bookId, Number((e.currentTarget as HTMLSelectElement).value));
  }
</script>

<header class="navbar bg-base-200 border-b border-base-300 px-3 gap-2">
  <div class="flex flex-wrap items-center gap-2">
    <select class="select select-sm select-bordered" value={study.bookId} onchange={onBookChange}>
      {#each books as book}
        <option value={book.id}>{book.name}</option>
      {/each}
    </select>
    <select class="select select-sm select-bordered w-24" value={String(study.chapter)} onchange={onChapterChange}>
      {#each Array.from({ length: books.find((b) => b.id === study.bookId)?.chapters ?? 1 }, (_, i) => i + 1) as ch}
        <option value={String(ch)}>{ch}</option>
      {/each}
    </select>
  </div>
  <div class="flex-1"></div>
  <button type="button" class="btn btn-sm btn-ghost" onclick={() => (study.settingsOpen = !study.settingsOpen)}>Settings</button>
</header>
