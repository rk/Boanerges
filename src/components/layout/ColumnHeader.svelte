<script lang="ts">
  import BookOpen from "@lucide/svelte/icons/book-open";
  import Columns2 from "@lucide/svelte/icons/columns-2";
  import Languages from "@lucide/svelte/icons/languages";
  import Link2 from "@lucide/svelte/icons/link-2";
  import List from "@lucide/svelte/icons/list";
  import Menu from "@lucide/svelte/icons/menu";
  import Pencil from "@lucide/svelte/icons/pencil";
  import Search from "@lucide/svelte/icons/search";
  import StickyNote from "@lucide/svelte/icons/sticky-note";
  import type { Component } from "svelte";
  import type { Snippet } from "svelte";
  import { COLUMN_CONTENT_LABELS } from "../../lib/columns/catalog";
  import { setColumnContent, study } from "../../lib/study.svelte.ts";
  import type { ColumnType } from "../../lib/study.svelte.ts";
  import { availableColumnOptions } from "../../lib/studyLayout";

  let {
    contentType,
    slotIndex,
    showViewSelector = false,
    children,
  }: {
    contentType: ColumnType | "primary-bible";
    slotIndex?: number;
    showViewSelector?: boolean;
    children?: Snippet;
  } = $props();

  const icons: Partial<Record<ColumnType | "primary-bible", Component>> = {
    notes: StickyNote,
    scribe: Pencil,
    search: Search,
    "cross-references": Link2,
    comparison: Columns2,
    "verse-list": List,
    dictionary: BookOpen,
    "bible-secondary": Languages,
  };

  let menuOpen = $state(false);
  let triggerEl = $state<HTMLButtonElement | null>(null);
  let panelEl = $state<HTMLUListElement | null>(null);

  const columnOptions = $derived(
    slotIndex === undefined
      ? []
      : availableColumnOptions(slotIndex, study.columns, study.translationBId, study.translationCId),
  );

  const title = $derived(
    contentType === "primary-bible" ? "" : (COLUMN_CONTENT_LABELS[contentType] ?? contentType),
  );

  const Icon = $derived(icons[contentType]);

  function closeMenu(): void {
    menuOpen = false;
  }

  function selectView(type: ColumnType): void {
    if (slotIndex === undefined) {
      return;
    }
    setColumnContent(slotIndex, type);
    closeMenu();
  }

  function toggleMenu(event: MouseEvent): void {
    event.stopPropagation();
    menuOpen = !menuOpen;
  }

  $effect(() => {
    if (!menuOpen) {
      return;
    }

    function handlePointerDown(event: PointerEvent): void {
      const target = event.target as Node;
      if (triggerEl?.contains(target) || panelEl?.contains(target)) {
        return;
      }
      closeMenu();
    }

    document.addEventListener("pointerdown", handlePointerDown);
    return () => document.removeEventListener("pointerdown", handlePointerDown);
  });
</script>

<div class="bg-base-200/50 flex min-h-11 shrink-0 items-center gap-2 border-b border-base-300/60 px-3 py-2">
  {#if showViewSelector && slotIndex !== undefined}
    <div class="dropdown dropdown-start" class:dropdown-open={menuOpen} class:z-50={menuOpen}>
      <button
        bind:this={triggerEl}
        type="button"
        class="btn btn-ghost btn-sm btn-square"
        aria-label="Change column view"
        aria-expanded={menuOpen}
        aria-haspopup="menu"
        data-testid="column-view-menu"
        onclick={toggleMenu}
      >
        <Menu size={16} aria-hidden="true" />
      </button>
      <ul
        bind:this={panelEl}
        class="dropdown-content menu bg-base-100 rounded-box z-50 w-44 border border-base-300 p-1 shadow-lg"
        role="menu"
      >
        {#each columnOptions as option (option)}
          <li role="none">
            <button
              type="button"
              role="menuitem"
              class:menu-active={contentType === option}
              onclick={() => selectView(option)}
            >
              {COLUMN_CONTENT_LABELS[option]}
            </button>
          </li>
        {/each}
      </ul>
    </div>
  {/if}

  {#if contentType !== "primary-bible" && Icon}
    <Icon size={16} class="text-base-content/70 shrink-0" aria-hidden="true" />
    <span class="text-sm font-medium shrink-0">{title}</span>
  {/if}
  {#if children}
    <div class="min-w-0 flex-1">{@render children()}</div>
  {/if}
</div>
