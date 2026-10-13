<script lang="ts">
  import Bookmark from "@lucide/svelte/icons/bookmark";
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import FormattedVerseText from "../reader/FormattedVerseText.svelte";
  import { bible, fetchBooksForTranslations } from "../../lib/bible.svelte.ts";
  import { comparison, scheduleComparisonLookup } from "../../lib/comparison.svelte.ts";
  import { parseScriptureReference } from "../../lib/scriptureReference";
  import { study } from "../../lib/study.svelte.ts";
  import { activeBibleTranslationIds } from "../../lib/studyLayout";

  let { slotIndex = 0 }: { slotIndex?: number; scrollRefs?: [HTMLElement | null, HTMLElement | null]; onscroll?: () => void } =
    $props();

  let referenceInput = $state("");

  const parsedReference = $derived(parseScriptureReference(referenceInput, bible.books));
  const referenceError = $derived(
    referenceInput.trim() !== "" && parsedReference === null ? "Enter a reference like Mark 1:1" : null,
  );

  $effect(() => {
    void fetchBooksForTranslations(activeBibleTranslationIds(study));
  });

  $effect(() => {
    if (comparison.activeReference !== null) {
      referenceInput = comparison.activeReference;
      comparison.activeReference = null;
    }
  });

  $effect(() => {
    scheduleComparisonLookup(referenceInput, bible.books);
  });
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" data-testid="comparison-column">
  <ColumnHeader contentType="comparison" {slotIndex} showViewSelector>
    <label class="input input-bordered input-sm min-w-0 flex-1">
      <Bookmark size={14} aria-hidden="true" />
      <input
        type="search"
        class="grow py-2"
        placeholder="Mark 1:1"
        data-testid="comparison-reference"
        bind:value={referenceInput}
        aria-label="Verse reference for comparison"
      />
    </label>
  </ColumnHeader>

  <div class="min-h-0 flex-1 overflow-y-auto px-2 py-2">
    {#if referenceError}
      <p class="text-error px-2 py-4 text-sm">{referenceError}</p>
    {:else if comparison.loading}
      <div class="flex justify-center py-8">
        <span class="loading loading-spinner loading-md text-primary"></span>
      </div>
    {:else if referenceInput.trim() === ""}
      <p class="px-2 py-4 text-sm opacity-60">Enter a verse reference to compare across translations.</p>
    {:else if comparison.rows.length === 0}
      <p class="px-2 py-4 text-sm opacity-60">No installed translations found.</p>
    {:else}
      <ul class="divide-y divide-base-300">
        {#each comparison.rows as row (row.translationId)}
          <li class="px-2 py-3" data-testid="comparison-row">
            <div class="mb-1 flex items-center gap-2">
              <span class="font-medium">{row.name}</span>
              <span class="badge badge-ghost badge-sm">{row.abbrev}</span>
            </div>
            <p class="text-sm">
              <FormattedVerseText text={row.text} />
            </p>
          </li>
        {/each}
      </ul>
    {/if}
  </div>
</div>
