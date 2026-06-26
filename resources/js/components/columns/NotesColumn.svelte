<script lang="ts">
    import CircleCheck from '@lucide/svelte/icons/circle-check';
    import LoaderCircle from '@lucide/svelte/icons/loader-circle';
    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import { bible, bookAbbrev } from '@/lib/bible.svelte.ts';
    import { fetchNotesDraft, notes, scheduleNotesSave } from '@/lib/notes.svelte.ts';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { study } from '@/lib/study.svelte.ts';

    let {
        slotIndex,
    }: {
        slotIndex: number;
    } = $props();

    let content = $state('');
    let loading = $state(true);

    const readerStyle = $derived(getReaderStyle());
    const currentBook = $derived(bible.books.find((book) => book.id === study.bookId));

    $effect(() => {
        const bookId = study.bookId;
        const chapterNumber = study.chapter;
        let cancelled = false;

        loading = true;
        content = '';

        void fetchNotesDraft(bookId, chapterNumber)
            .then((draft) => {
                if (! cancelled) {
                    content = draft;
                    loading = false;
                }
            })
            .catch(() => {
                if (! cancelled) {
                    loading = false;
                }
            });

        return () => {
            cancelled = true;
        };
    });

    function handleInput(event: Event): void {
        content = (event.currentTarget as HTMLTextAreaElement).value;
        scheduleNotesSave(study.bookId, study.chapter, content);
    }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col" style={readerStyle}>
    <ColumnHeader contentType="notes" {slotIndex} showViewSelector>
        {#snippet children()}
            <ChapterHeading
                title="{currentBook?.name ?? bookAbbrev(study.bookId)} {study.chapter}"
            />
            {#if notes.saveStatus === 'saving'}
                <LoaderCircle size={14} class="text-base-content/60 shrink-0 animate-spin" aria-label="Saving" />
            {:else if notes.saveStatus === 'saved'}
                <CircleCheck size={14} class="text-success shrink-0" aria-label="Saved" />
            {/if}
        {/snippet}
    </ColumnHeader>

    {#if loading}
        <div class="flex flex-1 items-center justify-center p-8">
            <span class="loading loading-spinner loading-lg text-primary"></span>
        </div>
    {:else}
        <textarea
            class="textarea textarea-ghost min-h-0 flex-1 resize-none rounded-none px-6 py-4 text-base leading-relaxed focus:outline-none"
            value={content}
            oninput={handleInput}
            placeholder="Notes for this chapter…"
            aria-label="Chapter notes"
        ></textarea>
    {/if}
</div>
