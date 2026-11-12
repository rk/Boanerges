<script lang="ts">
  import type { Chapter } from "../../mygo";
  import { BibleService } from "../../mygo";
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import ParagraphText from "../reader/ParagraphText.svelte";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { bible, fetchBooksForTranslations } from "../../lib/bible.svelte.ts";
  import { schedulePersistStudy, setSecondaryTranslation, study } from "../../lib/study.svelte.ts";

  let {
    translationId,
    primary = false,
    slotIndex = 0,
    scrollRef = $bindable<HTMLElement | null>(null),
    registerSlotScroll,
    syncLinkedScrollFromPrimary,
    onscroll,
  }: {
    translationId: string;
    primary?: boolean;
    slotIndex?: number;
    columnType?: string;
    scrollRef?: HTMLElement | null;
    /** Parent-owned stable callback — avoids $effect loop when parent passes fresh inline fns each render. */
    registerSlotScroll?: (slotIndex: number, el: HTMLElement | null) => void;
    /** Called after secondary scroll pane registers so primary position can propagate (brief desync window otherwise). */
    syncLinkedScrollFromPrimary?: () => void;
    onscroll?: (source: HTMLElement) => void;
  } = $props();

  let scrollEl = $state<HTMLElement | null>(null);

  $effect(() => {
    if (primary) {
      scrollRef = scrollEl;
    }
  });

  /** Mount-only registration — $effect + parent setState caused effect_update_depth_exceeded. */
  function registerSecondaryScrollPane(node: HTMLElement) {
    if (!primary) {
      registerSlotScroll?.(slotIndex, node);
      queueMicrotask(() => {
        if (study.scrollSync) {
          syncLinkedScrollFromPrimary?.();
        }
      });
    }
    return {
      destroy() {
        if (!primary) {
          registerSlotScroll?.(slotIndex, null);
        }
      },
    };
  }

  let chapter = $state<Chapter | null>(null);
  let chapterLoc = $state<{ bookId: string; chapter: number } | null>(null);
  let loading = $state(true);

  const translationOptions = $derived(bible.translations.map((t) => ({ id: t.id, abbrev: t.abbrev })));

  $effect(() => {
    const bookId = study.bookId;
    const chapterNumber = study.chapter;
    loading = true;
    chapter = null;
    chapterLoc = null;
    let cancelled = false;
    void BibleService.chapter(translationId, bookId, chapterNumber).then((resp) => {
      if (cancelled) {
        return;
      }
      chapter = resp.chapter;
      chapterLoc = { bookId, chapter: chapterNumber };
      loading = false;
    });
    return () => {
      cancelled = true;
    };
  });

  function onSecondaryTranslationChange(e: Event) {
    setSecondaryTranslation(slotIndex, (e.currentTarget as HTMLSelectElement).value);
  }

  function onPrimaryTranslationChange(e: Event) {
    study.translationId = (e.currentTarget as HTMLSelectElement).value;
    void fetchBooksForTranslations([study.translationId]);
    schedulePersistStudy();
  }

</script>

<div class="flex h-full min-h-0 flex-col">
  {#if primary}
    <ColumnHeader contentType="primary-bible">
      <select
        class="select select-bordered select-xs w-full max-w-full"
        data-testid="translation-select"
        value={translationId}
        onchange={onPrimaryTranslationChange}
      >
        {#each translationOptions as t (t.id)}
          <option value={t.id}>{t.abbrev}</option>
        {/each}
      </select>
    </ColumnHeader>
  {:else}
    <ColumnHeader contentType="bible-secondary" slotIndex={slotIndex} showViewSelector>
      <select
        class="select select-bordered select-xs w-full max-w-full"
        data-testid="secondary-translation-select"
        value={translationId}
        onchange={onSecondaryTranslationChange}
      >
        {#each translationOptions as t (t.id)}
          <option value={t.id}>{t.abbrev}</option>
        {/each}
      </select>
    </ColumnHeader>
  {/if}
  <article
    bind:this={scrollEl}
    use:registerSecondaryScrollPane
    class="prose max-w-none min-h-0 flex-1 overflow-y-auto p-4"
    style={readerStyle()}
    data-testid={primary ? "bible-reader-primary" : "bible-reader-secondary"}
    onscroll={(e) => onscroll?.(e.currentTarget as HTMLElement)}
  >
    {#if loading}
      <p>Loading…</p>
    {:else if chapter}
      <h2 class="chapter-heading" data-testid="chapter-heading">{chapter.book} {chapter.chapter}</h2>
      <ParagraphText
        verses={chapter.verses}
        bookId={chapterLoc?.bookId ?? study.bookId}
        chapter={chapterLoc?.chapter ?? study.chapter}
        showContextMenu={primary}
      />
    {/if}
  </article>
</div>
