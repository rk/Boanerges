<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import VerseListCombobox from "../ui/VerseListCombobox.svelte";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import { StudyData } from "../../mygo";
  import { bible, fetchBooksForTranslations } from "../../lib/bible.svelte.ts";
  import { formatScriptureReference, parseScriptureReference } from "../../lib/scriptureReference";
  import {
    navigateToVerse,
    schedulePersistStudy,
    study,
    verseListEntryKey,
  } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let lists = $state<{ id: string; title: string }[]>([]);
  let active = $state<{ id: string; title: string; entries: { bookId: string; chapter: number; verse: number }[] } | null>(
    null,
  );
  let titleDraft = $state("");
  let error = $state<string | null>(null);
  let entryDrafts = $state<string[]>([]);

  async function refresh() {
    lists = await StudyData.verseListIndex();
  }

  async function loadActive(id: string | null) {
    if (!id) {
      active = null;
      titleDraft = "";
      return;
    }
    const list = await StudyData.verseListGet(id);
    active = { id: list.id, title: list.title, entries: list.entries };
    titleDraft = list.title;
    entryDrafts = list.entries.map((e) => formatScriptureReference(e.bookId, e.chapter, e.verse));
  }

  $effect(() => {
    study.verseListTick;
    study.verseListActiveId;
    void (async () => {
      await refresh();
      await loadActive(study.verseListActiveId);
    })();
  });

  async function saveList(entries: { bookId: string; chapter: number; verse: number }[]) {
    error = null;
    const title = (titleDraft || active?.title || "My list").trim();
    if (!title) {
      error = "Title required";
      return;
    }
    try {
      const saved = await StudyData.verseListSave(title, entries, active?.id ?? study.verseListActiveId ?? "");
      study.verseListActiveId = saved.id;
      study.verseListTick += 1;
      schedulePersistStudy();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function ensureActiveList(): Promise<void> {
    if (active?.id) {
      return;
    }
    const existing = lists.find((l) => l.title === "My list");
    if (existing) {
      study.verseListActiveId = existing.id;
      await loadActive(existing.id);
      return;
    }
    const saved = await StudyData.verseListSave("My list", [], "");
    study.verseListActiveId = saved.id;
    study.verseListTick += 1;
    await refresh();
    await loadActive(saved.id);
  }

  async function addCurrentVerse() {
    if (study.selectedVerse === null) {
      return;
    }
    await ensureActiveList();
    const entries = [...(active?.entries ?? [])];
    const next = { bookId: study.bookId, chapter: study.chapter, verse: study.selectedVerse };
    if (entries.some((e) => verseListEntryKey(e) === verseListEntryKey(next))) {
      return;
    }
    entries.push(next);
    await saveList(entries);
  }

  async function createList(title?: string) {
    const name = (title?.trim() || "New list").trim();
    const saved = await StudyData.verseListSave(name, [], "");
    study.verseListActiveId = saved.id;
    study.verseListTick += 1;
    schedulePersistStudy();
    await refresh();
    await loadActive(saved.id);
  }

  async function deleteActiveList() {
    if (!active?.id) {
      return;
    }
    await StudyData.verseListDelete(active.id);
    study.verseListActiveId = null;
    active = null;
    titleDraft = "";
    study.verseListTick += 1;
    schedulePersistStudy();
  }

  function selectList(id: string) {
    study.verseListActiveId = id;
    schedulePersistStudy();
    void loadActive(id);
  }

  const comboboxOptions = $derived(
    lists.map((list) => ({
      id: list.id,
      title: list.title,
      entryCount:
        study.verseListActiveId === list.id && active ? active.entries.length : undefined,
    })),
  );

  function goToEntry(entry: { bookId: string; chapter: number; verse: number }) {
    navigateToVerse(entry.bookId, entry.chapter, entry.verse);
  }

  async function booksForParse(): Promise<void> {
    if (bible.books.length === 0) {
      await fetchBooksForTranslations([study.translationId]);
    }
  }

  async function commitEntryEdit(index: number) {
    error = null;
    const text = entryDrafts[index]?.trim() ?? "";
    if (!text) {
      return;
    }
    await booksForParse();
    const parsed = parseScriptureReference(text, bible.books);
    if (!parsed) {
      error = `Could not parse “${text}”`;
      if (active) {
        entryDrafts[index] = formatScriptureReference(
          active.entries[index].bookId,
          active.entries[index].chapter,
          active.entries[index].verse,
        );
      }
      return;
    }
    const entries = [...(active?.entries ?? [])];
    const key = verseListEntryKey(parsed);
    if (entries.some((e, i) => i !== index && verseListEntryKey(e) === key)) {
      error = "That reference is already in the list";
      return;
    }
    entries[index] = parsed;
    await saveList(entries);
  }

  async function removeEntry(index: number) {
    if (!active) {
      return;
    }
    const entries = active.entries.filter((_, i) => i !== index);
    await saveList(entries);
  }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" data-testid="verse-list-column">
  <ColumnHeader contentType="verse-list" {slotIndex} showViewSelector />
  <div class="flex min-h-0 min-w-0 flex-1 flex-col gap-0 border-t border-base-300">
    <section class="shrink-0 border-b border-base-300 px-3 py-2" aria-label="Verse lists">
      <div class="mb-2 flex flex-wrap items-center gap-2">
        <div class="min-w-0 flex-1 basis-48">
          <VerseListCombobox
            options={comboboxOptions}
            activeId={study.verseListActiveId}
            bind:draft={titleDraft}
            onSelect={selectList}
            onCreate={(title) => void createList(title)}
            onTitleCommit={() => {
              if (active) {
                void saveList(active.entries);
              }
            }}
          />
        </div>
        <button type="button" class="btn btn-xs" onclick={addCurrentVerse} disabled={study.selectedVerse === null}>
          Add verse
        </button>
        {#if active}
          <button type="button" class="btn btn-xs btn-ghost" data-testid="verse-list-delete" onclick={deleteActiveList}>
            Delete
          </button>
        {/if}
      </div>
      {#if error}
        <p class="text-error mb-1 text-xs">{error}</p>
      {/if}
    </section>

    <section class="flex min-h-0 min-w-0 flex-1 flex-col px-3 py-2" aria-label="Verses in list">
      {#if !active}
        <p class="text-sm opacity-60">Select or create a list, or add a verse from the Bible column.</p>
      {:else if active.entries.length === 0}
        <p class="text-sm opacity-60">No verses yet.</p>
      {:else}
        <ul class="flex min-h-0 w-full min-w-0 flex-1 flex-col gap-1 overflow-y-auto p-0">
          {#each active.entries as entry, index (verseListEntryKey(entry) + "-" + index)}
            <li class="flex w-full min-w-0 items-center gap-1 border-b border-base-200 py-1 last:border-0">
              <input
                class="input input-bordered input-xs min-w-0 flex-1 font-medium"
                bind:value={entryDrafts[index]}
                aria-label="Verse reference"
                onkeydown={(e) => {
                  if (e.key === "Enter") {
                    void commitEntryEdit(index);
                  }
                }}
                onblur={() => void commitEntryEdit(index)}
              />
              <button
                type="button"
                class="btn btn-xs btn-square btn-ghost shrink-0"
                aria-label="Go to verse"
                onclick={() => {
                  const e = active?.entries[index];
                  if (e) {
                    goToEntry(e);
                  }
                }}
              >
                →
              </button>
              <button
                type="button"
                class="btn btn-xs btn-square btn-ghost text-error shrink-0"
                aria-label="Remove verse"
                onclick={() => void removeEntry(index)}
              >
                <Trash2 size={14} aria-hidden="true" />
              </button>
            </li>
          {/each}
        </ul>
      {/if}
    </section>
  </div>
</div>
