<script lang="ts">
    import CircleCheck from '@lucide/svelte/icons/circle-check';
    import EllipsisVertical from '@lucide/svelte/icons/ellipsis-vertical';
    import LoaderCircle from '@lucide/svelte/icons/loader-circle';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import ScribeEditor from '@/components/scribe/ScribeEditor.svelte';
    import { fetchChapter, peekChapter } from '@/lib/bible.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { scribe } from '@/lib/scribe.svelte.ts';
    import { study } from '@/lib/study.svelte.ts';
    import type { Chapter } from '@/lib/types/bible';

    type ScribeEditorHandle = {
        openPreview(): void;
        setLineBreaksFromPrimary(): void;
        resetLineBreaks(): void;
    };

    let {
        slotIndex,
    }: {
        slotIndex: number;
    } = $props();

    let currentChapter = $state<Chapter | null>(null);
    let loading = $state(true);
    let lastLocationKey: string | null = null;
    let editor = $state<ScribeEditorHandle | null>(null);
    let actionsMenuOpen = $state(false);
    let actionsTriggerEl = $state<HTMLButtonElement | null>(null);
    let actionsPanelEl = $state<HTMLUListElement | null>(null);

    const readerStyle = $derived(getReaderStyle());
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

            if (
                actionsTriggerEl?.contains(target) ||
                actionsPanelEl?.contains(target)
            ) {
                return;
            }

            closeActionsMenu();
        }

        document.addEventListener('pointerdown', handlePointerDown);

        return () => {
            document.removeEventListener('pointerdown', handlePointerDown);
        };
    });

    $effect(() => {
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        const translationId = study.translationId;
        const locationKey = `${bookId}:${chapterNumber}`;
        const navigated =
            lastLocationKey !== null && locationKey !== lastLocationKey;
        const isInitialLoad = lastLocationKey === null;
        let cancelled = false;

        lastLocationKey = locationKey;

        if (isInitialLoad || navigated) {
            loading = true;
            currentChapter = null;
        } else {
            const cached = peekChapter(translationId, bookId, chapterNumber);

            if (cached) {
                currentChapter = cached;
                loading = false;
            }
        }

        void fetchChapter(translationId, bookId, chapterNumber)
            .then((chapter) => {
                if (!cancelled) {
                    currentChapter = chapter;
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

<div class="flex h-full min-h-0 min-w-0 flex-col" style={readerStyle}>
    <ColumnHeader contentType="scribe" {slotIndex} showViewSelector>
        {#if currentChapter}
            <ChapterHeading
                title="{currentChapter.book} {currentChapter.chapter}"
                compact={true}
            />
        {/if}
        {#if scribe.saveStatus === 'saving'}
            <LoaderCircle
                size={14}
                class="text-base-content/60 shrink-0 animate-spin"
                aria-label="Saving"
            />
        {:else if scribe.saveStatus === 'saved'}
            <CircleCheck
                size={14}
                class="text-success shrink-0"
                aria-label="Saved"
            />
        {/if}
        <div
            class="dropdown dropdown-end ml-auto"
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
    </ColumnHeader>

    {#if loading || !currentChapter}
        <div class="flex flex-1 items-center justify-center p-8">
            <span class="loading loading-spinner loading-lg text-primary"
            ></span>
        </div>
    {:else}
        <ScribeEditor
            bind:this={editor}
            book={currentChapter.book}
            chapter={currentChapter.chapter}
            verses={currentChapter.verses}
        />
    {/if}
</div>
