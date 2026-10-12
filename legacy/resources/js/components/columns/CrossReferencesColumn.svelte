<script lang="ts">
    import Bookmark from '@lucide/svelte/icons/bookmark';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import { bible } from '@/lib/bible.svelte.ts';
    import {
        crossrefs,
        scheduleCrossReferenceLookup,
    } from '@/lib/crossrefs.svelte.ts';
    import {
        formatCrossReference,
        parseScriptureReference,
    } from '@/lib/scriptureReference';
    import { goToVerseReference } from '@/lib/study.svelte.ts';
    import type { CrossReference } from '@/lib/types/bible';

    let {
        slotIndex,
    }: {
        slotIndex: number;
    } = $props();

    let referenceInput = $state('');

    const parsedReference = $derived(
        parseScriptureReference(referenceInput, bible.books),
    );
    const referenceError = $derived(
        referenceInput.trim() !== '' && parsedReference === null
            ? 'Enter a reference like Mark 1:1'
            : null,
    );

    $effect(() => {
        if (crossrefs.activeReference !== null) {
            referenceInput = crossrefs.activeReference;
            crossrefs.activeReference = null;
        }
    });

    $effect(() => {
        scheduleCrossReferenceLookup(referenceInput, bible.books);
    });

    function goTo(ref: CrossReference): void {
        goToVerseReference(ref.bookId, ref.chapter, ref.verse, ref.endVerse);
    }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col">
    <ColumnHeader contentType="cross-references" {slotIndex} showViewSelector>
        <label class="input input-bordered input-sm min-w-0 flex-1">
            <Bookmark size={14} aria-hidden="true" />
            <input
                type="search"
                class="grow py-2"
                placeholder="Mark 1:1"
                bind:value={referenceInput}
                aria-label="Verse reference for cross references"
            />
        </label>
    </ColumnHeader>

    <div class="min-h-0 flex-1 overflow-y-auto px-2 py-2">
        {#if referenceError}
            <p class="text-error px-2 py-4 text-sm">{referenceError}</p>
        {:else if crossrefs.loading}
            <div class="flex justify-center py-8">
                <span class="loading loading-spinner loading-md text-primary"
                ></span>
            </div>
        {:else if referenceInput.trim() === ''}
            <p class="text-base-content/60 px-2 py-4 text-sm">
                Enter a verse reference to look up cross references.
            </p>
        {:else if crossrefs.references.length === 0}
            <p class="text-base-content/60 px-2 py-4 text-sm">
                No cross references found.
            </p>
        {:else}
            <ul class="divide-base-300 divide-y">
                {#each crossrefs.references as ref (ref.bookId + ref.chapter + ref.verse + ref.rank)}
                    <li>
                        <button
                            type="button"
                            class="hover:bg-base-200 flex w-full items-center justify-between px-2 py-3 text-left"
                            onclick={() => goTo(ref)}
                        >
                            <span class="font-medium">
                                {formatCrossReference(ref, bible.books)}
                            </span>
                            <span class="badge badge-ghost badge-sm"
                                >{ref.rank}</span
                            >
                        </button>
                    </li>
                {/each}
            </ul>
        {/if}
    </div>
</div>
