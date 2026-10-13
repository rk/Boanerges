<script lang="ts">
  import ArrowLeft from "@lucide/svelte/icons/arrow-left";
  import { SvelteMap } from "svelte/reactivity";
  import type { Book } from "../../mygo";
  import { bible, fetchBooksForTranslations, getBooksForTranslation } from "../../lib/bible.svelte.ts";
  import { CANON_NT_BOOK_IDS, CANON_OT_BOOK_IDS, canonBookName } from "../../lib/canonBookNames";
  import { setBook, setChapter, study } from "../../lib/study.svelte.ts";
  import { activeBibleTranslationIds, isBookAvailableInTranslations } from "../../lib/studyLayout";

  type Step = "book" | "chapter";

  type BookRow = {
    id: string;
    name: string;
    chapters: number;
    available: boolean;
  };

  let {
    step = $bindable("book" as Step),
    onDone,
  }: {
    step?: Step;
    onDone?: () => void;
  } = $props();

  let pendingBookId = $state<string | null>(null);

  const primaryBooks = $derived(getBooksForTranslation(study.translationId) ?? bible.books);
  const activeTranslationIds = $derived(activeBibleTranslationIds(study));

  const booksByTranslationMap = $derived.by(() => {
    const map = new SvelteMap<string, readonly Book[]>();
    for (const translationId of activeTranslationIds) {
      const books = getBooksForTranslation(translationId);
      if (books) {
        map.set(translationId, books);
      }
    }
    return map;
  });

  const booksReady = $derived(
    activeTranslationIds.every((translationId) => getBooksForTranslation(translationId) !== undefined),
  );

  function buildBookRows(ids: readonly string[]): BookRow[] {
    return ids.map((id) => {
      const primaryBook = primaryBooks.find((book) => book.id === id);
      return {
        id,
        name: primaryBook?.name ?? canonBookName(id),
        chapters: primaryBook?.chapters ?? 1,
        available:
          booksReady && isBookAvailableInTranslations(id, activeTranslationIds, booksByTranslationMap),
      };
    });
  }

  const otBooks = $derived(buildBookRows(CANON_OT_BOOK_IDS));
  const ntBooks = $derived(buildBookRows(CANON_NT_BOOK_IDS));

  const pendingBook = $derived(
    primaryBooks.find((book) => book.id === pendingBookId) ??
      otBooks.find((book) => book.id === pendingBookId) ??
      ntBooks.find((book) => book.id === pendingBookId),
  );

  const chapterOptions = $derived(Array.from({ length: pendingBook?.chapters ?? 1 }, (_, index) => index + 1));

  function selectBook(bookId: string): void {
    pendingBookId = bookId;
    step = "chapter";
  }

  function selectChapter(chapterNumber: number): void {
    if (pendingBookId && pendingBookId !== study.bookId) {
      setBook(pendingBookId);
    }
    setChapter(chapterNumber);
    onDone?.();
  }

  function goBack(): void {
    step = "book";
    pendingBookId = null;
  }

  $effect(() => {
    void fetchBooksForTranslations(activeTranslationIds);
  });
</script>

<div class="bg-base-100 text-base-content min-h-0 p-3" data-testid="book-chapter-picker-panel">
  {#if step === "chapter"}
    <button type="button" class="btn btn-ghost btn-xs mb-2 gap-1 px-1" onclick={goBack}>
      <ArrowLeft size={14} aria-hidden="true" />
      Back
    </button>
  {/if}

  {#if step === "book"}
    {#if !booksReady || bible.booksLoading}
      <div class="flex items-center justify-center py-8">
        <span class="loading loading-spinner loading-md text-primary"></span>
      </div>
    {:else}
      <div class="max-h-[380px] overflow-x-hidden overflow-y-auto overscroll-y-contain">
        <div class="grid grid-cols-2 gap-3">
          <section class="min-w-0">
            <p class="menu-title px-0">Old Testament</p>
            <div class="columns-2 w-full space-y-1">
              {#each otBooks as book (book.id)}
                <button
                  type="button"
                  disabled={!book.available}
                  class="btn btn-info btn-xs block w-full min-w-0 truncate"
                  class:btn-active={study.bookId === book.id}
                  class:opacity-40={!book.available}
                  onclick={() => selectBook(book.id)}
                >
                  {book.name}
                </button>
              {/each}
            </div>
          </section>
          <section class="min-w-0">
            <p class="menu-title px-0">New Testament</p>
            <div class="columns-2 w-full space-y-1">
              {#each ntBooks as book (book.id)}
                <button
                  type="button"
                  disabled={!book.available}
                  class="btn btn-primary btn-xs block w-full min-w-0 truncate"
                  class:btn-active={study.bookId === book.id}
                  class:opacity-40={!book.available}
                  onclick={() => selectBook(book.id)}
                >
                  {book.name}
                </button>
              {/each}
            </div>
          </section>
        </div>
      </div>
    {/if}
  {:else}
    <p class="menu-title px-0">{pendingBook?.name ?? "Chapter"}</p>
    <div class="grid max-h-64 grid-cols-6 gap-1 overflow-y-auto">
      {#each chapterOptions as chapterNumber (chapterNumber)}
        <button
          type="button"
          class="btn btn-xs"
          class:btn-primary={pendingBookId === study.bookId && study.chapter === chapterNumber}
          onclick={() => selectChapter(chapterNumber)}
        >
          {chapterNumber}
        </button>
      {/each}
    </div>
  {/if}
</div>
