<script lang="ts">
  import Pilcrow from "@lucide/svelte/icons/pilcrow";
  import FormattedVerseText from "./FormattedVerseText.svelte";
  import type { Verse } from "../../mygo";
  import { AppChrome } from "../../mygo";
  import {
    addVerseToList,
    copyVerseReference,
    openComparisonForVerse,
    openCrossReferencesForVerse,
    openDictionaryForWord,
    setActiveVerse,
    study,
  } from "../../lib/study.svelte.ts";
  import { clampMenuPosition, wordFromVerseContext } from "../../lib/verseWord.ts";

  const useNativeFlyout = import.meta.env.VITE_E2E !== "1";

  let {
    verse,
    bookId,
    chapter,
    showContextMenu = false,
  }: {
    verse: Verse;
    bookId: string;
    chapter: number;
    showContextMenu?: boolean;
  } = $props();

  let menuOpen = $state(false);
  let menuX = $state(0);
  let menuY = $state(0);
  let menuWord = $state("");

  const highlighted = $derived(
    showContextMenu &&
      study.crossRefsTarget !== null &&
      study.crossRefsTarget.bookId === bookId &&
      study.crossRefsTarget.chapter === chapter &&
      study.crossRefsTarget.verse === verse.number,
  );

  async function onContextMenu(event: MouseEvent): Promise<void> {
    if (!showContextMenu) {
      return;
    }
    event.preventDefault();
    const word = wordFromVerseContext(event, verse.text);

    if (useNativeFlyout) {
      await AppChrome.openVerseActionsFlyout(
        Math.round(event.clientX),
        Math.round(event.clientY),
        bookId,
        chapter,
        verse.number,
        word,
      );
      return;
    }

    const pos = clampMenuPosition(event.clientX, event.clientY);
    menuX = pos.x;
    menuY = pos.y;
    menuWord = word;
    menuOpen = true;
  }

  function closeMenu(): void {
    menuOpen = false;
  }

  async function run(action: () => void | Promise<void>): Promise<void> {
    await action();
    closeMenu();
  }

  $effect(() => {
    if (!menuOpen) {
      return;
    }
    function onPointerDown(event: PointerEvent): void {
      const target = event.target as HTMLElement;
      if (target.closest("[data-verse-menu]")) {
        return;
      }
      closeMenu();
    }
    document.addEventListener("pointerdown", onPointerDown);
    return () => document.removeEventListener("pointerdown", onPointerDown);
  });
</script>

<p class:mt-4={verse.paragraphStart} data-verse={verse.number} oncontextmenu={onContextMenu}>
  {#if verse.paragraphStart}
    <Pilcrow size={14} class="mr-1 inline opacity-70" aria-hidden="true" />
  {/if}
  <span
    id="verse-{verse.number}"
    class={highlighted ? "rounded-sm bg-primary/15 px-0.5 [box-decoration-break:clone]" : undefined}
    role={showContextMenu ? "button" : undefined}
    tabindex={showContextMenu ? -1 : undefined}
    onclick={() => showContextMenu && setActiveVerse(verse.number)}
  >
    <sup class="text-primary mr-1 text-xs">{verse.number}</sup>
    <FormattedVerseText text={verse.text} />
  </span>
</p>

{#if menuOpen}
  <ul
    data-verse-menu
    class="menu bg-base-100 rounded-box fixed z-[2000] w-52 border border-base-300 p-1 shadow-lg"
    style="top: {menuY}px; left: {menuX}px;"
    role="menu"
  >
    <li>
      <button
        type="button"
        role="menuitem"
        onclick={() => run(() => openComparisonForVerse(bookId, chapter, verse.number))}
      >
        Compare
      </button>
    </li>
    <li>
      <button
        type="button"
        role="menuitem"
        onclick={() => run(() => openCrossReferencesForVerse(bookId, chapter, verse.number))}
      >
        Cross-references
      </button>
    </li>
    <li>
      <button
        type="button"
        role="menuitem"
        onclick={() => run(() => addVerseToList(bookId, chapter, verse.number))}
      >
        Add to verse list
      </button>
    </li>
    <li>
      <button
        type="button"
        role="menuitem"
        onclick={() => run(() => copyVerseReference(bookId, chapter, verse.number))}
      >
        Copy reference
      </button>
    </li>
    {#if menuWord.length >= 2}
      <li>
        <button type="button" role="menuitem" onclick={() => run(() => openDictionaryForWord(menuWord))}>
          Dictionary: {menuWord}
        </button>
      </li>
    {/if}
  </ul>
{/if}
