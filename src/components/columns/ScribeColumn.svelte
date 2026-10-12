<script lang="ts">
  import { StudyData } from "../../mygo";
  import { study } from "../../lib/study.svelte.ts";

  let draft = $state("");

  $effect(() => {
    void StudyData.scribeGet(study.bookId, study.chapter).then((verses) => {
      draft = verses.map((v) => `${v.verse}. ${v.text}`).join("\n");
    });
  });

  function save() {
    const verses = draft
      .split("\n")
      .map((line) => line.trim())
      .filter(Boolean)
      .map((line) => {
        const m = line.match(/^(\d+)\.\s*(.*)$/);
        return { verse: m ? Number(m[1]) : 0, text: m ? m[2] : line };
      })
      .filter((v) => v.verse > 0);
    void StudyData.scribePut(study.bookId, study.chapter, verses);
  }
</script>

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">Scribe</div>
  <textarea
    class="textarea textarea-ghost h-full resize-none rounded-none p-4 font-mono text-sm"
    data-testid="scribe-editor"
    bind:value={draft}
    onblur={save}
    placeholder="1. Verse draft…"
  ></textarea>
</div>
