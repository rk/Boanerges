<script lang="ts">
  import CircleCheck from "@lucide/svelte/icons/circle-check";
  import LoaderCircle from "@lucide/svelte/icons/loader-circle";
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import { StudyData } from "../../mygo";
  import { notes, scheduleNotesSave } from "../../lib/notes.svelte.ts";
  import { schedulePersistStudy, study } from "../../lib/study.svelte.ts";

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let content = $state("");
  let loadGen = 0;

  $effect(() => {
    const bookId = study.bookId;
    const chapter = study.chapter;
    const gen = ++loadGen;
    void StudyData.notesGet(bookId, chapter).then((text) => {
      if (gen === loadGen) {
        content = text;
      }
    });
  });

  function onInput() {
    scheduleNotesSave(study.bookId, study.chapter, content);
    schedulePersistStudy();
  }

  function onBlur() {
    scheduleNotesSave(study.bookId, study.chapter, content);
    schedulePersistStudy();
  }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" data-testid="notes-column">
  <ColumnHeader contentType="notes" {slotIndex} showViewSelector>
    <div class="ml-auto flex shrink-0 items-center gap-2">
      {#if notes.saveStatus === "saving"}
        <span class="text-base-content/60 flex items-center gap-1 text-xs">
          <LoaderCircle size={14} class="animate-spin" aria-hidden="true" />
          Saving…
        </span>
      {:else if notes.saveStatus === "saved"}
        <span class="text-success flex items-center gap-1 text-xs">
          <CircleCheck size={14} aria-hidden="true" />
          Saved
        </span>
      {/if}
    </div>
  </ColumnHeader>
  <textarea
    class="textarea textarea-ghost min-h-0 w-full min-w-0 flex-1 resize-none rounded-none border-0 p-4 focus:outline-none"
    data-testid="notes-editor"
    bind:value={content}
    oninput={onInput}
    onblur={onBlur}
    placeholder="Chapter notes…"
  ></textarea>
</div>
