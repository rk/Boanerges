<script lang="ts">
  import { StudyData } from "../../mygo";
  import { study } from "../../lib/study.svelte.ts";

  let lists = $state<{ id: string; title: string }[]>([]);
  let active = $state<{ title: string; entries: { bookId: string; chapter: number; verse: number }[] } | null>(null);

  async function refresh() {
    lists = await StudyData.verseListIndex();
    if (study.verseListActiveId) {
      active = await StudyData.verseListGet(study.verseListActiveId);
    }
  }

  $effect(() => {
    void refresh();
  });

  async function addCurrentVerse() {
    const entries = active?.entries ?? [];
    entries.push({ bookId: study.bookId, chapter: study.chapter, verse: study.selectedVerse });
    const saved = await StudyData.verseListSave(active?.title ?? "My list", entries, study.verseListActiveId ?? "");
    study.verseListActiveId = saved.id;
    await refresh();
  }
</script>

<div class="flex h-full flex-col" data-testid="verse-list-column">
  <div class="flex items-center justify-between border-b border-base-300 px-3 py-2 text-sm">
    <span class="font-medium">Verse list</span>
    <button type="button" class="btn btn-xs" onclick={addCurrentVerse}>Add verse</button>
  </div>
  <ul class="menu p-2">
    {#each lists as list}
      <li><button type="button" class={study.verseListActiveId === list.id ? "active" : ""} onclick={async () => { study.verseListActiveId = list.id; active = await StudyData.verseListGet(list.id); }}>{list.title}</button></li>
    {/each}
  </ul>
  {#if active}
    <ul class="flex-1 overflow-y-auto p-3 text-sm">
      {#each active.entries as entry}
        <li>{entry.bookId} {entry.chapter}:{entry.verse}</li>
      {/each}
    </ul>
  {/if}
</div>
