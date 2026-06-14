<script lang="ts">
    import { bible } from '@/lib/bible.svelte.ts';
    import { study, setBook } from '@/lib/study.svelte.ts';

    const oldTestament = $derived(bible.books.filter((book) => book.testament === 'ot'));
    const newTestament = $derived(bible.books.filter((book) => book.testament === 'nt'));
</script>

<div class="grid grid-cols-2 gap-2">
    <div>
        <p class="menu-title px-0 text-xs">Old Testament</p>
        <ul class="menu menu-xs rounded-box bg-base-100 max-h-40 w-full overflow-y-auto p-1">
            {#if bible.booksLoading}
                <li><span class="text-base-content/60 px-2 py-1 text-xs">Loading…</span></li>
            {:else}
                {#each oldTestament as book (book.id)}
                    <li>
                        <button
                            type="button"
                            class:menu-active={study.bookId === book.id}
                            onclick={() => setBook(book.id)}
                        >
                            {book.name}
                        </button>
                    </li>
                {/each}
            {/if}
        </ul>
    </div>

    <div>
        <p class="menu-title px-0 text-xs">New Testament</p>
        <ul class="menu menu-xs rounded-box bg-base-100 max-h-40 w-full overflow-y-auto p-1">
            {#if bible.booksLoading}
                <li><span class="text-base-content/60 px-2 py-1 text-xs">Loading…</span></li>
            {:else}
                {#each newTestament as book (book.id)}
                    <li>
                        <button
                            type="button"
                            class:menu-active={study.bookId === book.id}
                            onclick={() => setBook(book.id)}
                        >
                            {book.name}
                        </button>
                    </li>
                {/each}
            {/if}
        </ul>
    </div>
</div>
