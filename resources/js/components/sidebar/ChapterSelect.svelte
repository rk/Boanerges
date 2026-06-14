<script lang="ts">
    import { bible } from '@/lib/bible.svelte.ts';
    import { study, setChapter } from '@/lib/study.svelte.ts';

    const currentBook = $derived(bible.books.find((book) => book.id === study.bookId));
    const chapterOptions = $derived(
        Array.from({ length: currentBook?.chapters ?? 1 }, (_, index) => index + 1),
    );
</script>

<label class="form-control w-full">
    <div class="label py-1">
        <span class="label-text text-xs">Chapter</span>
    </div>
    <select
        class="select select-bordered select-sm w-full"
        value={study.chapter}
        onchange={(event) => setChapter(Number(event.currentTarget.value))}
        disabled={! currentBook}
    >
        {#each chapterOptions as chapterNumber (chapterNumber)}
            <option value={chapterNumber}>{chapterNumber}</option>
        {/each}
    </select>
</label>

<div class="mt-2 grid grid-cols-6 gap-1">
    {#each chapterOptions.slice(0, 30) as chapterNumber (chapterNumber)}
        <button
            type="button"
            class="btn btn-xs"
            class:btn-primary={study.chapter === chapterNumber}
            onclick={() => setChapter(chapterNumber)}
        >
            {chapterNumber}
        </button>
    {/each}
</div>
