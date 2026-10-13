<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { StudyData } from "../../mygo";
  import { schedulePersistStudy, setLocation, study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let lists = $state<{ id: string; title: string }[]>([]);
  let active = $state<{ id: string; title: string; entries: { bookId: string; chapter: number; verse: number }[] } | null>(
    null,
  );
  let titleDraft = $state("");
  let error = $state<string | null>(null);

  async function refresh() {
    lists = await StudyData.verseListIndex();
    if (study.verseListActiveId) {
      const list = await StudyData.verseListGet(study.verseListActiveId);
      active = { id: list.id, title: list.title, entries: list.entries };
      titleDraft = list.title;
    }
  }

  $effect(() => {
    void refresh();
  });

  async function saveList(entries: { bookId: string; chapter: number; verse: number }[]) {
    error = null;
    try {
      const saved = await StudyData.verseListSave(titleDraft || active?.title || "My list", entries, active?.id ?? "");
      study.verseListActiveId = saved.id;
      schedulePersistStudy();
      await refresh();
    } catch (err) {
      error = err instanceof Error ? err.message : String(err);
    }
  }

  async function addCurrentVerse() {
    const entries = [...(active?.entries ?? [])];
    entries.push({ bookId: study.bookId, chapter: study.chapter, verse: study.selectedVerse });
    await saveList(entries);
  }

  async function deleteActiveList() {
    if (!active?.id) return;
    await StudyData.verseListDelete(active.id);
    study.verseListActiveId = null;
    active = null;
    schedulePersistStudy();
    await refresh();
  }

  function goToEntry(entry: { bookId: string; chapter: number; verse: number }) {
    setLocation(entry.bookId, entry.chapter);
    study.selectedVerse = entry.verse;
  }
</script>

<div class="flex h-full flex-col" data-testid="verse-list-column">
  <ColumnHeader contentType="verse-list" {slotIndex} showViewSelector />
  <div class="flex flex-wrap items-center gap-2 border-b border-base-300 px-3 py-2">
    <input class="input input-bordered input-xs flex-1" bind:value={titleDraft} placeholder="List title" />
    <button type="button" class="btn btn-xs" onclick={addCurrentVerse}>Add verse</button>
    {#if active}
      <button type="button" class="btn btn-xs btn-ghost" data-testid="verse-list-delete" onclick={deleteActiveList}>Delete</button>
    {/if}
  </div>
  {#if error}
    <p class="text-error px-3 py-1 text-xs">{error}</p>
  {/if}
  <ul class="menu p-2">
    {#each lists as list (list.id)}
      <li>
        <button
          type="button"
          class={study.verseListActiveId === list.id ? "active" : ""}
          onclick={async () => {
            study.verseListActiveId = list.id;
            schedulePersistStudy();
            active = await StudyData.verseListGet(list.id);
            titleDraft = active.title;
          }}
        >
          {list.title}
        </button>
      </li>
    {/each}
  </ul>
  {#if active}
    <ul class="flex-1 overflow-y-auto p-3 text-sm">
      {#each active.entries as entry, index (index)}
        <li>
          <button type="button" class="link link-hover" onclick={() => goToEntry(entry)}>
            {entry.bookId} {entry.chapter}:{entry.verse}
          </button>
        </li>
      {/each}
    </ul>
  {/if}
</div>
