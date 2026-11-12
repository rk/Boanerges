<script lang="ts">
  import ChevronDown from "@lucide/svelte/icons/chevron-down";

  export type VerseListOption = { id: string; title: string; entryCount?: number };

  let {
    options,
    activeId = null,
    draft = $bindable(""),
    disabled = false,
    onSelect,
    onCreate,
    onTitleCommit,
  }: {
    options: VerseListOption[];
    activeId?: string | null;
    draft?: string;
    disabled?: boolean;
    onSelect: (id: string) => void;
    onCreate: (title: string) => void;
    onTitleCommit: () => void;
  } = $props();

  let open = $state(false);
  let highlight = $state(0);
  let searchQuery = $state("");
  let skipCommitOnBlur = $state(false);
  let rootEl = $state<HTMLDivElement | null>(null);
  let panelEl = $state<HTMLDivElement | null>(null);
  let titleInputEl = $state<HTMLInputElement | null>(null);
  let searchInputEl = $state<HTMLInputElement | null>(null);

  const filterText = $derived(searchQuery.trim().toLowerCase());

  const filtered = $derived.by(() => {
    const items = !filterText
      ? [...options]
      : options.filter((o) => o.title.toLowerCase().includes(filterText));
    items.sort((a, b) => {
      const al = a.title.toLowerCase();
      const bl = b.title.toLowerCase();
      const aStart = filterText ? (al.startsWith(filterText) ? 0 : 1) : 0;
      const bStart = filterText ? (bl.startsWith(filterText) ? 0 : 1) : 0;
      if (aStart !== bStart) {
        return aStart - bStart;
      }
      return al.localeCompare(bl);
    });
    return items;
  });

  const exactMatch = $derived(
    options.some((o) => o.title.trim().toLowerCase() === filterText && filterText.length > 0),
  );

  const showCreate = $derived(!exactMatch);

  const rowCount = $derived(filtered.length + (showCreate ? 1 : 0));

  function close(): void {
    open = false;
    highlight = 0;
    searchQuery = "";
  }

  function focusSearchSoon(): void {
    queueMicrotask(() => searchInputEl?.focus());
  }

  function openPanel(): void {
    if (disabled) {
      return;
    }
    searchQuery = "";
    open = true;
    const activeIndex = filtered.findIndex((o) => o.id === activeId);
    highlight = activeIndex >= 0 ? activeIndex : 0;
    focusSearchSoon();
  }

  function chooseList(id: string): void {
    skipCommitOnBlur = true;
    onSelect(id);
    close();
    titleInputEl?.blur();
  }

  function chooseCreate(): void {
    skipCommitOnBlur = true;
    const title = searchQuery.trim() || "New list";
    onCreate(title);
    close();
    titleInputEl?.blur();
  }

  function commitHighlighted(): void {
    if (!open || rowCount === 0) {
      close();
      return;
    }
    if (highlight < filtered.length) {
      chooseList(filtered[highlight].id);
      return;
    }
    if (showCreate) {
      chooseCreate();
    }
  }

  function onTitleKeydown(event: KeyboardEvent): void {
    if (event.key === "ArrowDown") {
      event.preventDefault();
      if (!open) {
        openPanel();
        return;
      }
      searchInputEl?.focus();
      return;
    }
    if (event.key === "Escape") {
      event.preventDefault();
      const active = options.find((o) => o.id === activeId);
      draft = active?.title ?? "";
      close();
      return;
    }
    if (event.key === "Enter") {
      event.preventDefault();
      onTitleCommit();
      close();
    }
  }

  function onSearchKeydown(event: KeyboardEvent): void {
    if (event.key === "ArrowDown") {
      event.preventDefault();
      highlight = Math.min(highlight + 1, rowCount - 1);
      return;
    }
    if (event.key === "ArrowUp") {
      event.preventDefault();
      if (highlight === 0) {
        return;
      }
      highlight = Math.max(highlight - 1, 0);
      return;
    }
    if (event.key === "Escape") {
      event.preventDefault();
      close();
      titleInputEl?.focus();
      return;
    }
    if (event.key === "Enter") {
      event.preventDefault();
      if (rowCount > 0) {
        commitHighlighted();
      }
    }
  }

  function toggleChevron(event: MouseEvent): void {
    event.preventDefault();
    if (open) {
      close();
      titleInputEl?.focus();
      return;
    }
    openPanel();
  }

  $effect(() => {
    if (!open) {
      return;
    }
    if (highlight >= rowCount) {
      highlight = Math.max(0, rowCount - 1);
    }

    function handlePointerDown(event: PointerEvent): void {
      const target = event.target as Node;
      if (rootEl?.contains(target)) {
        return;
      }
      close();
    }

    document.addEventListener("pointerdown", handlePointerDown);
    return () => document.removeEventListener("pointerdown", handlePointerDown);
  });

  $effect(() => {
    searchQuery;
    filterText;
    if (!open) {
      return;
    }
    if (highlight >= rowCount) {
      highlight = Math.max(0, rowCount - 1);
    } else if (filterText) {
      highlight = 0;
    }
  });
</script>

<div bind:this={rootEl} class="dropdown w-full min-w-0" class:dropdown-open={open} class:z-50={open}>
  <label class="input input-bordered input-xs flex min-w-0 w-full items-center gap-1">
    <input
      bind:this={titleInputEl}
      type="text"
      class="min-w-0 grow bg-transparent outline-none"
      aria-haspopup="listbox"
      aria-expanded={open}
      aria-controls="verse-list-combobox-panel"
      aria-label="List title"
      placeholder="List title"
      data-testid="verse-list-combobox"
      {disabled}
      bind:value={draft}
      onkeydown={onTitleKeydown}
      onblur={() => {
        if (skipCommitOnBlur) {
          skipCommitOnBlur = false;
          return;
        }
        if (!open) {
          onTitleCommit();
        }
      }}
    />
    <button
      type="button"
      class="btn btn-ghost btn-xs btn-square shrink-0"
      tabindex={-1}
      aria-label={open ? "Close list options" : "Open list options"}
      aria-expanded={open}
      {disabled}
      onmousedown={(e) => e.preventDefault()}
      onclick={toggleChevron}
    >
      <ChevronDown size={14} class="transition-transform {open ? 'rotate-180' : ''}" aria-hidden="true" />
    </button>
  </label>

  <div
    bind:this={panelEl}
    id="verse-list-combobox-panel"
    class="dropdown-content bg-base-100 rounded-box z-50 mt-1 w-full min-w-0 border border-base-300 p-0 shadow-lg"
  >
    <div class="border-base-300 border-b px-2 py-1">
      <input
        bind:this={searchInputEl}
        type="search"
        class="input input-bordered input-xs w-full"
        role="combobox"
        aria-autocomplete="list"
        aria-label="Search verse lists"
        aria-controls="verse-list-combobox-options"
        placeholder="Search lists"
        data-testid="verse-list-combobox-search"
        {disabled}
        bind:value={searchQuery}
        onkeydown={onSearchKeydown}
      />
    </div>
    <ul
      id="verse-list-combobox-options"
      class="menu max-h-40 w-full min-w-0 overflow-y-auto p-1 pt-0"
      role="listbox"
    >
      {#each filtered as option, index (option.id)}
        <li class="w-full min-w-0" role="presentation">
          <button
            type="button"
            role="option"
            aria-selected={activeId === option.id}
            class="w-full truncate text-left"
            class:active={activeId === option.id}
            class:menu-active={open && highlight === index}
            onmousedown={(e) => e.preventDefault()}
            onclick={() => chooseList(option.id)}
          >
            {option.title}
            {#if option.entryCount != null}
              <span class="badge badge-ghost badge-xs ml-1">{option.entryCount}</span>
            {/if}
          </button>
        </li>
      {:else}
        {#if !showCreate}
          <li class="w-full px-2 py-1 text-xs opacity-60" role="presentation">No lists yet</li>
        {/if}
      {/each}
      {#if showCreate}
        <li class="w-full min-w-0 border-t border-base-200" role="presentation">
          <button
            type="button"
            role="option"
            class="w-full truncate text-left"
            class:menu-active={open && highlight === filtered.length}
            onmousedown={(e) => e.preventDefault()}
            onclick={chooseCreate}
          >
            {#if searchQuery.trim()}
              Create “{searchQuery.trim()}”
            {:else}
              Create new list
            {/if}
          </button>
        </li>
      {/if}
    </ul>
  </div>
</div>
