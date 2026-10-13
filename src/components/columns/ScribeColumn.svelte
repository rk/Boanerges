<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { StudyData, type ScribeVerse } from "../../mygo";
  import { study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let draft = $state("");
  let loaded = $state<ScribeVerse[]>([]);

  $effect(() => {
    void StudyData.scribeGet(study.bookId, study.chapter).then((verses) => {
      loaded = verses;
      draft = verses.map((v) => `${v.verse}. ${v.text}`).join("\n");
    });
  });

  function save() {
    const byVerse = new Map(loaded.map((v) => [v.verse, v]));
    const verses = draft
      .split("\n")
      .map((line) => line.trim())
      .filter(Boolean)
      .map((line) => {
        const m = line.match(/^(\d+)\.\s*(.*)$/);
        const verse = m ? Number(m[1]) : 0;
        const text = m ? m[2] : line;
        const prev = byVerse.get(verse);
        return { verse, text, paragraphStart: prev?.paragraphStart ?? false };
      })
      .filter((v) => v.verse > 0);
    void StudyData.scribePut(study.bookId, study.chapter, verses);
    loaded = verses;
  }
</script>

<div class="flex h-full flex-col">
  <ColumnHeader contentType="scribe" {slotIndex} showViewSelector />
  <textarea
    class="textarea textarea-ghost h-full resize-none rounded-none p-4 font-mono text-sm"
    data-testid="scribe-editor"
    bind:value={draft}
    onblur={save}
    placeholder="1. Verse draft…"
  ></textarea>
</div>
