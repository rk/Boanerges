<script lang="ts">
    import FolderOpen from '@lucide/svelte/icons/folder-open';
    import Save from '@lucide/svelte/icons/save';
    import Trash2 from '@lucide/svelte/icons/trash-2';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import FormattedVerseText from '@/components/reader/FormattedVerseText.svelte';
    import { bible } from '@/lib/bible.svelte.ts';
    import { formatScriptureReference } from '@/lib/scriptureReference';
    import {
        goToVerseReference,
        setVerseListActiveId,
        setVerseListShowContentSetting,
        study,
    } from '@/lib/study.svelte.ts';
    import {
        loadSavedLists,
        loadSavedVerseList,
        loadVerseTexts,
        removeVerseFromList,
        saveVerseList,
        setVerseListShowContent,
        verseList,
        verseListNeedsTextLoad,
    } from '@/lib/verseList.svelte.ts';

    let {
        slotIndex,
    }: {
        slotIndex: number;
    } = $props();

    let loadMenuOpen = $state(false);

    $effect(() => {
        void loadSavedLists();
    });

    $effect(() => {
        if (
            !verseListNeedsTextLoad(
                verseList.entries,
                verseList.showContent,
                study.translationId,
                verseList.loading,
            )
        ) {
            return;
        }

        void loadVerseTexts(study.translationId);
    });

    async function handleSave(): Promise<void> {
        const saved = await saveVerseList();

        if (saved && verseList.savedListId) {
            setVerseListActiveId(verseList.savedListId);
        }
    }

    async function handleLoad(id: string): Promise<void> {
        if (
            verseList.entries.length > 0 &&
            !confirm('Replace the current verse list?')
        ) {
            return;
        }

        loadMenuOpen = false;
        await loadSavedVerseList(id);
        setVerseListActiveId(id);

        verseList.entries = verseList.entries.map((entry) => ({
            ...entry,
            label: formatScriptureReference(
                entry.bookId,
                entry.chapter,
                entry.verse,
                bible.books,
            ),
        }));
    }

    function toggleShowContent(event: Event): void {
        const checked = (event.currentTarget as HTMLInputElement).checked;
        setVerseListShowContent(checked);
        setVerseListShowContentSetting(checked);
    }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col">
    <ColumnHeader contentType="verse-list" {slotIndex} showViewSelector>
        <input
            type="text"
            class="input input-bordered input-sm min-w-0 flex-1"
            bind:value={verseList.title}
            aria-label="Verse list title"
            placeholder="Untitled list"
        />
        <button
            type="button"
            class="btn btn-ghost btn-sm btn-square"
            aria-label="Save verse list"
            disabled={verseList.saving}
            onclick={handleSave}
        >
            <Save size={16} aria-hidden="true" />
        </button>
        <div class="dropdown dropdown-end" class:dropdown-open={loadMenuOpen}>
            <button
                type="button"
                class="btn btn-ghost btn-sm btn-square"
                aria-label="Load verse list"
                aria-haspopup="menu"
                onclick={() => {
                    loadMenuOpen = !loadMenuOpen;
                }}
            >
                <FolderOpen size={16} aria-hidden="true" />
            </button>
            <ul
                class="dropdown-content menu bg-base-100 rounded-box border-base-300 z-10 w-52 border p-2 shadow"
                role="menu"
            >
                {#if verseList.savedLists.length === 0}
                    <li class="text-base-content/60 px-2 py-2 text-sm">
                        No saved lists
                    </li>
                {:else}
                    {#each verseList.savedLists as saved (saved.id)}
                        <li role="none">
                            <button
                                type="button"
                                role="menuitem"
                                onclick={() => handleLoad(saved.id)}
                            >
                                {saved.title}
                            </button>
                        </li>
                    {/each}
                {/if}
            </ul>
        </div>
        <label class="label cursor-pointer gap-2 py-0">
            <span class="label-text text-xs">Show text</span>
            <input
                type="checkbox"
                class="toggle toggle-sm"
                checked={verseList.showContent}
                onchange={toggleShowContent}
                aria-label="Show verse text"
            />
        </label>
    </ColumnHeader>

    <div class="min-h-0 flex-1 overflow-y-auto px-2 py-2">
        {#if verseList.loading && verseList.entries.length === 0}
            <div class="flex justify-center py-8">
                <span class="loading loading-spinner loading-md text-primary"
                ></span>
            </div>
        {:else if verseList.entries.length === 0}
            <p class="text-base-content/60 px-2 py-4 text-sm">
                Right-click a verse to add it to this list.
            </p>
        {:else}
            <ul class="divide-base-300 divide-y">
                {#each verseList.entries as entry, index (entry.bookId + entry.chapter + entry.verse)}
                    <li class="flex items-start gap-2 px-2 py-3">
                        <button
                            type="button"
                            class="hover:bg-base-200 min-w-0 flex-1 rounded px-1 text-left"
                            onclick={() =>
                                goToVerseReference(
                                    entry.bookId,
                                    entry.chapter,
                                    entry.verse,
                                )}
                        >
                            <span class="font-medium">{entry.label}</span>
                            {#if verseList.showContent && entry.text}
                                <p
                                    class="reader-prose text-base-content/90 mt-1 text-sm"
                                >
                                    <FormattedVerseText text={entry.text} />
                                </p>
                            {/if}
                        </button>
                        <button
                            type="button"
                            class="btn btn-ghost btn-sm btn-square shrink-0"
                            aria-label="Remove {entry.label}"
                            onclick={() => removeVerseFromList(index)}
                        >
                            <Trash2 size={16} aria-hidden="true" />
                        </button>
                    </li>
                {/each}
            </ul>
        {/if}
    </div>
</div>
