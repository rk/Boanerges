<script lang="ts">
    import { Search } from '@lucide/svelte';
    import { router } from '@inertiajs/svelte';
    import { search, scheduleSearch, cancelSearch } from '@/lib/search.svelte.ts';
    import { study, setBook, setChapter } from '@/lib/study.svelte.ts';
    import type { SearchResult } from '@/lib/types/bible';

    let dialog: HTMLDialogElement | undefined = $state();
    let query = $state('');

    $effect(() => {
        if (search.open) {
            dialog?.showModal();
        } else {
            dialog?.close();
        }
    });

    $effect(() => {
        const q = query.trim();
        const translationId = study.translationId;

        scheduleSearch(q, translationId);
    });

    function handleClose(): void {
        cancelSearch();
        search.open = false;
        query = '';
    }

    function goToResult(result: SearchResult): void {
        setBook(result.bookId);
        setChapter(result.chapter);
        handleClose();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box flex max-h-[85vh] max-w-lg flex-col">
        <h3 class="text-lg font-bold">Search Scripture</h3>

        <label class="input input-bordered mt-4">
            <Search size="14" />
            <input type="search" class="grow py-2" placeholder="Search…" bind:value={query} />
        </label>

        <div class="mt-4 min-h-0 flex-1 overflow-y-auto">
            {#if search.loading}
                <div class="flex justify-center py-8">
                    <span class="loading loading-spinner loading-md text-primary"></span>
                </div>
            {:else if query.trim().length < 2}
                <p class="text-base-content/60 py-4 text-sm">Type at least two characters.</p>
            {:else if search.results.length === 0}
                <p class="text-base-content/60 py-4 text-sm">No results found.</p>
            {:else}
                <ul class="divide-base-300 divide-y">
                    {#each search.results as result (result.bookId + result.chapter + result.verse)}
                        <li>
                            <button
                                type="button"
                                class="hover:bg-base-200 w-full px-2 py-3 text-left"
                                onclick={() => goToResult(result)}
                            >
                                <p class="text-sm font-medium uppercase">
                                    {result.bookId} {result.chapter}:{result.verse}
                                </p>
                                <p class="text-base-content/70 mt-1 text-sm">{@html result.snippet}</p>
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
        <button aria-label="Close search">close</button>
    </form>
</dialog>
