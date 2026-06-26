<script lang="ts">
    import NotesColumn from '@/components/columns/NotesColumn.svelte';
    import SearchColumn from '@/components/columns/SearchColumn.svelte';
    import BibleColumn from '@/components/columns/BibleColumn.svelte';
    import CrossReferencesColumn from '@/components/columns/CrossReferencesColumn.svelte';
    import ScribeColumn from '@/components/columns/ScribeColumn.svelte';
    import ChapterNavRail from '@/components/layout/ChapterNavRail.svelte';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { bibleColumnCount, secondaryTranslationForSlot } from '@/lib/studyLayout';
    import {
        setTranslation,
        setTranslationB,
        setTranslationC,
        study,
    } from '@/lib/study.svelte.ts';
    import type { ColumnContentType } from '@/lib/types/study';
    import { createVerseHighlightScroller } from '@/lib/verseHighlightScroll';

    type Slot = {
        type: ColumnContentType;
        slotIndex: number;
    };

    let primaryScroll = $state<HTMLElement | null>(null);
    let secondaryScrolls = $state<[HTMLElement | null, HTMLElement | null]>([null, null]);
    let syncing = $state(false);
    let suppressScrollSync = $state(false);

    const highlightScroller = createVerseHighlightScroller();
    const readerStyle = $derived(getReaderStyle());
    const showLeftNavRail = $derived(study.columnCount > 1);

    const slots = $derived(
        study.columns.map((type, slotIndex): Slot => ({ type, slotIndex })),
    );

    const bibleScrollRefs = $derived.by(() => {
        const refs: HTMLElement[] = [];

        if (primaryScroll) {
            refs.push(primaryScroll);
        }

        for (const slot of slots) {
            if (slot.type === 'bible-secondary') {
                const ref = secondaryScrolls[slot.slotIndex];

                if (ref) {
                    refs.push(ref);
                }
            }
        }

        return refs;
    });

    $effect(() => {
        if (! study.verseHighlight) {
            highlightScroller.reset();
        }
    });

    $effect(() => {
        const highlight = study.verseHighlight;
        const expectedRootCount = bibleColumnCount(study);

        if (! highlight || expectedRootCount === 0) {
            return;
        }

        let cancelled = false;

        void (async () => {
            suppressScrollSync = true;

            await highlightScroller.scrollTo(
                bibleScrollRefs,
                highlight,
                `${study.bookId}:${study.chapter}`,
                expectedRootCount,
                () => cancelled,
            );

            if (cancelled) {
                return;
            }

            await new Promise((resolve) => requestAnimationFrame(resolve));
            await new Promise((resolve) => requestAnimationFrame(resolve));
            suppressScrollSync = false;
        })();

        return () => {
            cancelled = true;
            suppressScrollSync = false;
        };
    });

    function syncScroll(source: HTMLElement, target: HTMLElement | null): void {
        if (! study.scrollSync || ! target || syncing || suppressScrollSync) {
            return;
        }

        syncing = true;

        const sourceRange = source.scrollHeight - source.clientHeight;
        const targetRange = target.scrollHeight - target.clientHeight;
        const ratio = sourceRange > 0 ? source.scrollTop / sourceRange : 0;

        target.scrollTop = ratio * targetRange;

        syncing = false;
    }

    function makeScrollHandler(slotIndex: number | 'primary'): () => void {
        return () => {
            const source = slotIndex === 'primary' ? primaryScroll : secondaryScrolls[slotIndex];

            if (! source) {
                return;
            }

            for (const ref of bibleScrollRefs) {
                if (ref !== source) {
                    syncScroll(source, ref);
                }
            }
        };
    }

    function secondaryTranslationChange(slotIndex: number): (id: string) => void {
        return slotIndex === 0 ? setTranslationB : setTranslationC;
    }
</script>

<div
    class="flex h-full min-h-0"
    class:mx-auto={study.columnCount === 1}
    class:max-w-prose={study.columnCount === 1}
    style={readerStyle}
>
    {#if showLeftNavRail}
        <ChapterNavRail />
    {/if}

    <div class="flex min-h-0 min-w-0 flex-1">
        <div class="min-w-0 flex-1">
            <BibleColumn
                translationId={study.translationId}
                showChapterNav={! showLeftNavRail}
                onTranslationChange={setTranslation}
                bind:scrollRef={primaryScroll}
                onscroll={makeScrollHandler('primary')}
            />
        </div>

        {#each slots as slot (slot.slotIndex)}
            <div class="min-w-0 flex-1">
                {#if slot.type === 'bible-secondary'}
                    <BibleColumn
                        translationId={secondaryTranslationForSlot(
                            slot.slotIndex,
                            study.translationBId,
                            study.translationCId,
                        )}
                        slotIndex={slot.slotIndex}
                        onTranslationChange={secondaryTranslationChange(slot.slotIndex)}
                        bind:scrollRef={secondaryScrolls[slot.slotIndex]}
                        onscroll={makeScrollHandler(slot.slotIndex)}
                    />
                {:else if slot.type === 'scribe'}
                    <ScribeColumn slotIndex={slot.slotIndex} />
                {:else if slot.type === 'notes'}
                    <NotesColumn slotIndex={slot.slotIndex} />
                {:else if slot.type === 'cross-references'}
                    <CrossReferencesColumn slotIndex={slot.slotIndex} />
                {:else}
                    <SearchColumn slotIndex={slot.slotIndex} />
                {/if}
            </div>
        {/each}
    </div>
</div>
