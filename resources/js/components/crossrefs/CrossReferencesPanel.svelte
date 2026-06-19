<script lang="ts">
    import { BookMarked } from '@lucide/svelte';
    import { crossrefs, loadCrossReferences } from '@/lib/crossrefs.svelte.ts';
    import { study, setBook, setChapter } from '@/lib/study.svelte.ts';
    import type { CrossReference } from '@/lib/types/bible';

    let dialog: HTMLDialogElement | undefined = $state();
    let selectedVerse = $state<number | null>(null);

    $effect(() => {
        if (crossrefs.open) {
            dialog?.showModal();

            if (selectedVerse !== null) {
                void loadCrossReferences(study.bookId, study.chapter, selectedVerse);
            }
        } else {
            dialog?.close();
        }
    });

    $effect(() => {
        if (crossrefs.open && selectedVerse !== null) {
            void loadCrossReferences(study.bookId, study.chapter, selectedVerse);
        }
    });

    export function openForVerse(verse: number): void {
        selectedVerse = verse;
        crossrefs.open = true;
    }

    function handleClose(): void {
        crossrefs.open = false;
        selectedVerse = null;
    }

    function goTo(ref: CrossReference): void {
        setBook(ref.bookId);
        setChapter(ref.chapter);
        handleClose();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box flex max-h-[85vh] max-w-lg flex-col">
        <h3 class="text-lg font-bold flex items-center gap-2">
            <BookMarked size="18" /> Cross References
        </h3>

        {#if selectedVerse !== null}
            <p class="text-base-content/70 mt-1 text-sm">
                {study.bookId.toUpperCase()} {study.chapter}:{selectedVerse}
            </p>
        {/if}

        <div class="mt-4 min-h-0 flex-1 overflow-y-auto">
            {#if crossrefs.loading}
                <div class="flex justify-center py-8">
                    <span class="loading loading-spinner loading-md text-primary"></span>
                </div>
            {:else if crossrefs.references.length === 0}
                <p class="text-base-content/60 py-4 text-sm">Select a verse or no references found.</p>
            {:else}
                <ul class="divide-base-300 divide-y">
                    {#each crossrefs.references as ref (ref.bookId + ref.chapter + ref.verse + ref.rank)}
                        <li>
                            <button
                                type="button"
                                class="hover:bg-base-200 flex w-full items-center justify-between px-2 py-3 text-left"
                                onclick={() => goTo(ref)}
                            >
                                <span class="font-medium uppercase">
                                    {ref.bookId} {ref.chapter}:{ref.verse}{#if ref.endVerse && ref.endVerse !== ref.verse}–{ref.endVerse}{/if}
                                </span>
                                <span class="badge badge-ghost badge-sm">{ref.rank}</span>
                            </button>
                        </li>
                    {/each}
                </ul>
            {/if}
        </div>

        <div class="modal-action">
            <form method="dialog">
                <button type="submit" class="btn">Close</button>
            </form>
        </div>
    </div>

    <form method="dialog" class="modal-backdrop">
        <button aria-label="Close cross references">close</button>
    </form>
</dialog>
