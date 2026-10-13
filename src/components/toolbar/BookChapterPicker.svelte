<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";
  import { AppChrome } from "../../mygo";
  import { bible, getBooksForTranslation } from "../../lib/bible.svelte.ts";
  import { study } from "../../lib/study.svelte.ts";
  import BookChapterPickerPanel from "./BookChapterPickerPanel.svelte";

  const useNativeFlyout = import.meta.env.VITE_E2E !== "1";

  let triggerEl = $state<HTMLButtonElement | null>(null);
  let panelEl = $state<HTMLDivElement | null>(null);
  let open = $state(false);
  let panelStyle = $state("");
  let step = $state<"book" | "chapter">("book");

  const primaryBooks = $derived(getBooksForTranslation(study.translationId) ?? bible.books);
  const currentBook = $derived(primaryBooks.find((book) => book.id === study.bookId));
  const locationLabel = $derived(currentBook ? `${currentBook.name} ${study.chapter}` : "Book & chapter");

  function updatePosition(): void {
    if (!triggerEl) {
      return;
    }
    const rect = triggerEl.getBoundingClientRect();
    panelStyle = `top:${rect.bottom + 4}px;left:${rect.left}px;`;
  }

  function closeFallback(): void {
    open = false;
    step = "book";
  }

  async function togglePicker(event: MouseEvent): Promise<void> {
    event.stopPropagation();
    if (useNativeFlyout && triggerEl) {
      const rect = triggerEl.getBoundingClientRect();
      await AppChrome.openBookChapterFlyout(
        Math.round(rect.x),
        Math.round(rect.y),
        Math.round(rect.width),
        Math.round(rect.height),
      );
      return;
    }
    open = !open;
    if (open) {
      step = "book";
      updatePosition();
    }
  }

  function portal(node: HTMLElement): { destroy: () => void } {
    document.body.appendChild(node);
    return {
      destroy() {
        node.remove();
      },
    };
  }

  $effect(() => {
    if (!open || useNativeFlyout) {
      return;
    }
    updatePosition();

    function handlePointerDown(event: PointerEvent): void {
      const target = event.target as Node;
      if (triggerEl?.contains(target) || panelEl?.contains(target)) {
        return;
      }
      closeFallback();
    }

    function handleReposition(): void {
      updatePosition();
    }

    document.addEventListener("pointerdown", handlePointerDown);
    window.addEventListener("resize", handleReposition);
    window.addEventListener("scroll", handleReposition, true);

    return () => {
      document.removeEventListener("pointerdown", handlePointerDown);
      window.removeEventListener("resize", handleReposition);
      window.removeEventListener("scroll", handleReposition, true);
    };
  });
</script>

<button
  bind:this={triggerEl}
  type="button"
  class="btn btn-ghost btn-sm gap-1"
  data-testid="book-chapter-picker"
  aria-label="Select book and chapter"
  aria-expanded={open}
  onclick={togglePicker}
>
  <span class="max-w-40 truncate">{locationLabel}</span>
  <ChevronDown size={14} aria-hidden="true" />
</button>

{#if open && !useNativeFlyout}
  <div
    use:portal
    bind:this={panelEl}
    class="bg-base-100 rounded-box fixed z-[1000] w-[min(36rem,calc(100vw-1.5rem))] overflow-x-hidden border border-base-300 shadow-lg"
    style={panelStyle}
    role="dialog"
    aria-label="Book and chapter picker"
  >
    <BookChapterPickerPanel bind:step onDone={closeFallback} />
  </div>
{/if}
