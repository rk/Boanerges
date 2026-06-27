<script lang="ts">
    import Search from '@lucide/svelte/icons/search';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import { canonBookName } from '@/lib/canonBookNames';
    import {
        loadMoreSearchResults,
        scheduleSearch,
        search,
    } from '@/lib/search.svelte.ts';
    import { goToVerseReference, study } from '@/lib/study.svelte.ts';
    import type { SearchResult } from '@/lib/types/bible';

    let {
        slotIndex,
    }: {
        slotIndex: number;
    } = $props();

    let query = $state('');

    $effect(() => {
        scheduleSearch(query, study.translationId);
    });

    function openResult(result: SearchResult): void {
        goToVerseReference(result.bookId, result.chapter, result.verse);
    }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col">
    <ColumnHeader contentType="search" {slotIndex} showViewSelector>
        <label class="input input-bordered input-sm min-w-0 flex-1">
            <Search size={14} aria-hidden="true" />
            <input
                type="search"
                class="grow py-2"
                placeholder="Search scripture…"
                bind:value={query}
            />
        </label>
    </ColumnHeader>

    <div class="min-h-0 flex-1 overflow-y-auto px-2 py-2">
        {#if search.loading && search.results.length === 0}
            <div class="flex justify-center py-8">
                <span class="loading loading-spinner loading-md text-primary"></span>
            </div>
        {:else if query.trim().length < 2}
            <p class="text-base-content/60 px-2 py-4 text-sm">Type at least two characters.</p>
        {:else if search.results.length === 0}
            <p class="text-base-content/60 px-2 py-4 text-sm">No results found.</p>
        {:else}
            <ul class="divide-base-300 divide-y">
                {#each search.results as result (result.bookId + result.chapter + result.verse)}
                    <li>
                        <button
                            type="button"
                            class="hover:bg-base-200 w-full rounded-lg px-2 py-3 text-left"
                            onclick={() => openResult(result)}
                        >
                            <p class="text-sm font-medium">
                                {canonBookName(result.bookId)} {result.chapter}:{result.verse}
                            </p>
                            <p class="text-base-content/70 mt-1 text-sm">{@html result.snippet}</p>
                        </button>
                    </li>
                {/each}
            </ul>

            {#if search.hasMore}
                <div class="flex justify-center py-4">
                    <button
                        type="button"
                        class="btn btn-sm btn-ghost"
                        disabled={search.loading}
                        onclick={() => loadMoreSearchResults(study.translationId)}
                    >
                        {#if search.loading}
                            <span class="loading loading-spinner loading-xs"></span>
                        {:else}
                            Load more
                        {/if}
                    </button>
                </div>
            {/if}
        {/if}
    </div>
</div>
