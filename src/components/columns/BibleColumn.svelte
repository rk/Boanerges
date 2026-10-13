<script lang="ts">
  import { TranslationService } from "../../mygo";
  import type { Chapter } from "../../mygo";
  import { BibleService } from "../../mygo";
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import VerseLine from "../reader/VerseLine.svelte";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { setSecondaryTranslation, study } from "../../lib/study.svelte.ts";

  let {
    translationId,
    primary = false,
    slotIndex = 0,
    scrollRef = $bindable<HTMLElement | null>(null),
    onscroll,
  }: {
    translationId: string;
    primary?: boolean;
    slotIndex?: number;
    columnType?: string;
    scrollRef?: HTMLElement | null;
    onscroll?: () => void;
  } = $props();

  let chapter = $state<Chapter | null>(null);
  let loading = $state(true);
  let translations = $state<{ id: string; abbrev: string }[]>([]);

  $effect(() => {
    if (!primary) {
      void TranslationService.installed().then((resp) => {
        translations = resp.translations.map((t) => ({ id: t.id, abbrev: t.abbrev }));
      });
    }
  });

  $effect(() => {
    loading = true;
    void BibleService.chapter(translationId, study.bookId, study.chapter).then((resp) => {
      chapter = resp.chapter;
      loading = false;
    });
  });

  function onSecondaryTranslationChange(e: Event) {
    setSecondaryTranslation(slotIndex, (e.currentTarget as HTMLSelectElement).value);
  }
</script>

<div class="flex h-full flex-col">
  {#if primary}
    <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">Primary · {translationId.toUpperCase()}</div>
  {:else}
    <ColumnHeader contentType="bible-secondary" slotIndex={slotIndex} showViewSelector>
      <select
        class="select select-bordered select-xs w-full max-w-full"
        data-testid="secondary-translation-select"
        value={translationId}
        onchange={onSecondaryTranslationChange}
      >
        {#each translations as t (t.id)}
          <option value={t.id}>{t.abbrev}</option>
        {/each}
      </select>
    </ColumnHeader>
  {/if}
  <article
    bind:this={scrollRef}
    class="prose max-w-none flex-1 overflow-y-auto p-4"
    style={readerStyle()}
    data-testid={primary ? "bible-reader-primary" : "bible-reader-secondary"}
    onscroll={onscroll}
  >
    {#if loading}
      <p>Loading…</p>
    {:else if chapter}
      <h2 data-testid="chapter-heading">{chapter.book} {chapter.chapter}</h2>
      {#each chapter.verses as verse (verse.number)}
        <VerseLine {verse} highlighted={primary && verse.number === study.selectedVerse} />
      {/each}
    {/if}
  </article>
</div>
