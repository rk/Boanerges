<script lang="ts">
  import { StudyData } from "../../mygo";
  import { schedulePersistStudy, study } from "../../lib/study.svelte.ts";

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

<div class="flex h-full flex-col">
  <div class="border-b border-base-300 px-3 py-2 text-sm font-medium">Notes</div>
  <textarea
    class="textarea textarea-ghost h-full resize-none rounded-none p-4"
    data-testid="notes-editor"
    bind:value={content}
    oninput={save}
    placeholder="Chapter notes…"
  ></textarea>
</div>
