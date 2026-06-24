<script lang="ts">
    import { BookMarked, Search } from '@lucide/svelte';
    import { bible } from '@/lib/bible.svelte.ts';
    import {
        closeCrossReferences,
        crossrefs,
        scheduleCrossReferenceLookup,
        scriptureReferenceForCurrentVerse,
    } from '@/lib/crossrefs.svelte.ts';
    import {
        formatScriptureReference,
        formatCrossReference,
        parseScriptureReference,
    } from '@/lib/scriptureReference';
    import { study, goToVerseReference } from '@/lib/study.svelte.ts';
    import type { CrossReference } from '@/lib/types/bible';

    let dialog: HTMLDialogElement | undefined = $state();
    let referenceInput = $state('');

    const parsedReference = $derived(parseScriptureReference(referenceInput, bible.books));
    const referenceLabel = $derived(
        parsedReference
            ? formatScriptureReference(
                parsedReference.bookId,
                parsedReference.chapter,
                parsedReference.verse,
                bible.books,
            )
            : null,
    );
    const referenceError = $derived(
        referenceInput.trim() !== '' && parsedReference === null
            ? 'Enter a reference like Mark 1:1'
            : null,
    );

    $effect(() => {
        if (crossrefs.open) {
            dialog?.showModal();

            if (crossrefs.initialReference !== null) {
                referenceInput = crossrefs.initialReference;
                crossrefs.initialReference = null;
            } else if (referenceInput.trim() === '') {
                referenceInput = scriptureReferenceForCurrentVerse(bible.books);
            }
        } else {
            dialog?.close();
        }
    });

    $effect(() => {
        if (! crossrefs.open) {
            return;
        }

        scheduleCrossReferenceLookup(referenceInput, bible.books);
    });

    export function openForVerse(verse: number): void {
        crossrefs.initialReference = formatScriptureReference(
            study.bookId,
            study.chapter,
            verse,
            bible.books,
        );
        crossrefs.open = true;
    }

    function handleClose(): void {
        referenceInput = '';
        closeCrossReferences();
    }

    function goTo(ref: CrossReference): void {
        goToVerseReference(ref.bookId, ref.chapter, ref.verse, ref.endVerse);
        handleClose();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box flex max-h-[85vh] max-w-lg flex-col">
        <h3 class="text-lg font-bold flex items-center gap-2">
            <BookMarked size="18" /> Cross References
        </h3>

        <label class="input input-bordered mt-4">
            <Search size="14" />
            <input
                type="search"
                class="grow py-2"
                placeholder="Mark 1:1"
                bind:value={referenceInput}
            />
        </label>

        {#if referenceError}
            <p class="text-error mt-2 text-sm">{referenceError}</p>
        {:else if referenceLabel}
            <p class="text-base-content/70 mt-2 text-sm">{referenceLabel}</p>
        {/if}

        <div class="mt-4 min-h-0 flex-1 overflow-y-auto">
            {#if crossrefs.loading}
                <div class="flex justify-center py-8">
                    <span class="loading loading-spinner loading-md text-primary"></span>
                </div>
            {:else if referenceInput.trim() === ''}
                <p class="text-base-content/60 py-4 text-sm">Enter a verse reference to look up cross references.</p>
            {:else if referenceError}
                <p class="text-base-content/60 py-4 text-sm">Could not parse that reference.</p>
            {:else if crossrefs.references.length === 0}
                <p class="text-base-content/60 py-4 text-sm">No cross references found.</p>
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
