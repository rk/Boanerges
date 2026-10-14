<script lang="ts">
  import CircleCheck from "@lucide/svelte/icons/circle-check";
  import EllipsisVertical from "@lucide/svelte/icons/ellipsis-vertical";
  import LoaderCircle from "@lucide/svelte/icons/loader-circle";
  import ColumnHeader from "../layout/ColumnHeader.svelte";
  import ScribeEditor from "../scribe/ScribeEditor.svelte";
  import { BibleService, type Chapter } from "../../mygo";
  import { readerStyle } from "../../lib/readability.svelte.ts";
  import { scribe } from "../../lib/scribe.svelte.ts";
  import { study } from "../../lib/study.svelte.ts";

  type ScribeEditorHandle = {
    openPreview(): void;
    setLineBreaksFromPrimary(): void;
    resetLineBreaks(): void;
  };

  let { slotIndex = 0 }: { slotIndex?: number; columnType?: string } = $props();

  let currentChapter = $state<Chapter | null>(null);
  let loading = $state(true);
  let editor = $state<ScribeEditorHandle | null>(null);
  let actionsMenuOpen = $state(false);
  let actionsTriggerEl = $state<HTMLButtonElement | null>(null);
  let actionsPanelEl = $state<HTMLUListElement | null>(null);

  const actionsDisabled = $derived(loading || !currentChapter);

  function closeActionsMenu(): void {
    actionsMenuOpen = false;
  }

  function toggleActionsMenu(event: MouseEvent): void {
    event.stopPropagation();
    actionsMenuOpen = !actionsMenuOpen;
  }

  $effect(() => {
    if (!actionsMenuOpen) {
      return;
    }

    function handlePointerDown(event: PointerEvent): void {
      const target = event.target as Node;

      if (actionsTriggerEl?.contains(target) || actionsPanelEl?.contains(target)) {
        return;
      }

      closeActionsMenu();
    }

    document.addEventListener("pointerdown", handlePointerDown);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
    };
  });

  $effect(() => {
    const bookId = study.bookId;
    const chapterNumber = study.chapter;
    const translationId = study.translationId;
    let cancelled = false;

    loading = true;
    currentChapter = null;

    void BibleService.chapter(translationId, bookId, chapterNumber)
      .then((resp) => {
        if (!cancelled) {
          currentChapter = resp.chapter;
          loading = false;
        }
      })
      .catch(() => {
        if (!cancelled) {
          loading = false;
        }
      });

    return () => {
      cancelled = true;
    };
  });
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" style={readerStyle()}>
  <ColumnHeader contentType="scribe" {slotIndex} showViewSelector>
    <div class="flex min-w-0 flex-1 items-center gap-2">
      {#if currentChapter}
        <span class="truncate text-sm opacity-80">{currentChapter.book} {currentChapter.chapter}</span>
      {/if}
      <div class="ml-auto flex shrink-0 items-center gap-2">
        {#if scribe.saveStatus === "saving"}
          <span class="text-base-content/60 flex items-center gap-1 text-xs">
            <LoaderCircle size={14} class="animate-spin" aria-hidden="true" />
            Saving…
          </span>
        {:else if scribe.saveStatus === "saved"}
          <span class="text-success flex items-center gap-1 text-xs">
            <CircleCheck size={14} aria-hidden="true" />
            Saved
          </span>
        {/if}
        <div
          class="dropdown dropdown-end"
          class:dropdown-open={actionsMenuOpen}
          class:z-50={actionsMenuOpen}
        >
          <button
            bind:this={actionsTriggerEl}
            type="button"
            class="btn btn-ghost btn-sm btn-square"
            aria-label="Scribe actions"
            aria-expanded={actionsMenuOpen}
            aria-haspopup="menu"
            disabled={actionsDisabled}
            onclick={toggleActionsMenu}
          >
            <EllipsisVertical size={16} aria-hidden="true" />
          </button>
          <ul
            bind:this={actionsPanelEl}
            class="dropdown-content menu bg-base-100 rounded-box z-50 w-56 border border-base-300 p-1 shadow-lg"
            role="menu"
          >
            <li role="none">
              <button
                type="button"
                role="menuitem"
                onclick={() => {
                  editor?.setLineBreaksFromPrimary();
                  closeActionsMenu();
                }}
              >
                Set Line Breaks from Primary
              </button>
            </li>
            <li role="none">
              <button
                type="button"
                role="menuitem"
                onclick={() => {
                  editor?.resetLineBreaks();
                  closeActionsMenu();
                }}
              >
                Reset Line Breaks
              </button>
            </li>
            <li role="none">
              <button
                type="button"
                role="menuitem"
                onclick={() => {
                  editor?.openPreview();
                  closeActionsMenu();
                }}
              >
                Preview
              </button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </ColumnHeader>

  {#if loading || !currentChapter}
    <div class="flex flex-1 items-center justify-center p-8">
      <span class="loading loading-lg loading-spinner text-primary"></span>
    </div>
  {:else}
    <ScribeEditor bind:this={editor} book={currentChapter.book} chapter={currentChapter.chapter} verses={currentChapter.verses} />
  {/if}
</div>
