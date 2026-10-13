<script lang="ts">
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { StudyData } from "../../mygo";
  import { schedulePersistStudy, study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let content = $state("");
  let timer: ReturnType<typeof setTimeout> | null = null;

  $effect(() => {
    void StudyData.notesGet(study.bookId, study.chapter).then((text) => {
      content = text;
    });
  });

  function save() {
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => {
      void StudyData.notesPut(study.bookId, study.chapter, content);
      schedulePersistStudy();
    }, 300);
  }
</script>

<div class="flex h-full min-h-0 flex-col">
  <ColumnHeader contentType="notes" {slotIndex} showViewSelector />
  <textarea
    class="textarea textarea-ghost min-h-0 flex-1 resize-none rounded-none p-4"
    data-testid="notes-editor"
    bind:value={content}
    oninput={save}
    placeholder="Chapter notes…"
  ></textarea>
</div>
