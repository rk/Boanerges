<script lang="ts">
    import ArrowLeft from '@lucide/svelte/icons/arrow-left';
    import ChevronDown from '@lucide/svelte/icons/chevron-down';
    import { SvelteMap } from 'svelte/reactivity';

    import {
        bible,
        fetchBooksForTranslations,
        getBooksForTranslation,
    } from '@/lib/bible.svelte.ts';
    import {
        CANON_NT_BOOK_IDS,
        CANON_OT_BOOK_IDS,
        canonBookName,
    } from '@/lib/canonBookNames';
    import { study, setBook, setChapter } from '@/lib/study.svelte.ts';
    import {
        activeBibleTranslationIds,
        isBookAvailableInTranslations,
    } from '@/lib/studyLayout';
    import type { Book } from '@/lib/types/bible';

    type Step = 'book' | 'chapter';

    type BookRow = {
        id: string;
        name: string;
        chapters: number;
        available: boolean;
    };

    let triggerEl = $state<HTMLButtonElement | null>(null);
    let panelEl = $state<HTMLDivElement | null>(null);
    let open = $state(false);
    let panelStyle = $state('');
    let step = $state<Step>('book');
    let pendingBookId = $state<string | null>(null);

    const primaryBooks = $derived(
        getBooksForTranslation(study.translationId) ?? bible.books,
    );

    const currentBook = $derived(
        primaryBooks.find((book) => book.id === study.bookId),
    );
    const locationLabel = $derived(
        currentBook ? `${currentBook.name} ${study.chapter}` : 'Book & chapter',
    );

    const activeTranslationIds = $derived(activeBibleTranslationIds(study));

    const booksByTranslationMap = $derived.by(() => {
        const map = new SvelteMap<string, readonly Book[]>();

        for (const translationId of activeTranslationIds) {
            const books = getBooksForTranslation(translationId);

            if (books) {
                map.set(translationId, books);
            }
        }

        return map;
    });

    const booksReady = $derived(
        activeTranslationIds.every(
            (translationId) =>
                getBooksForTranslation(translationId) !== undefined,
        ),
    );

    function buildBookRows(ids: readonly string[]): BookRow[] {
        return ids.map((id) => {
            const primaryBook = primaryBooks.find((book) => book.id === id);

            return {
                id,
                name: primaryBook?.name ?? canonBookName(id),
                chapters: primaryBook?.chapters ?? 1,
                available:
                    booksReady &&
                    isBookAvailableInTranslations(
                        id,
                        activeTranslationIds,
                        booksByTranslationMap,
                    ),
            };
        });
    }

    const otBooks = $derived(buildBookRows(CANON_OT_BOOK_IDS));
    const ntBooks = $derived(buildBookRows(CANON_NT_BOOK_IDS));

    const pendingBook = $derived(
        primaryBooks.find((book) => book.id === pendingBookId) ??
            otBooks.find((book) => book.id === pendingBookId) ??
            ntBooks.find((book) => book.id === pendingBookId) ??
            currentBook,
    );

    const chapterOptions = $derived(
        Array.from(
            { length: pendingBook?.chapters ?? 1 },
            (_, index) => index + 1,
        ),
    );

    function resetPicker(): void {
        step = 'book';
        pendingBookId = null;
    }

    function closePicker(): void {
        open = false;
        resetPicker();
    }

    function updatePosition(): void {
        if (!triggerEl) {
            return;
        }

        const rect = triggerEl.getBoundingClientRect();
        panelStyle = `top:${rect.bottom + 4}px;left:${rect.left}px;`;
    }

    function togglePicker(event: MouseEvent): void {
        event.stopPropagation();
        open = !open;

        if (open) {
            resetPicker();
            updatePosition();
        }
    }

    function selectBook(bookId: string): void {
        pendingBookId = bookId;
        step = 'chapter';
    }

    function selectChapter(chapterNumber: number): void {
        if (pendingBookId && pendingBookId !== study.bookId) {
            setBook(pendingBookId);
        }

        setChapter(chapterNumber);
        closePicker();
    }

    function goBack(): void {
        step = 'book';
        pendingBookId = null;
    }

    function portal(node: HTMLElement): { destroy: () => void } {
        document.body.appendChild(node);

        return {
            destroy() {
                node.remove();
            },
        };
    }

    $effect(() => {
        if (!open) {
            return;
        }

        void fetchBooksForTranslations(activeTranslationIds);
    });

    $effect(() => {
        if (!open) {
            return;
        }

        updatePosition();

        function handlePointerDown(event: PointerEvent): void {
            const target = event.target as Node;

            if (triggerEl?.contains(target) || panelEl?.contains(target)) {
                return;
            }

            closePicker();
        }

        function handleReposition(): void {
            updatePosition();
        }

        document.addEventListener('pointerdown', handlePointerDown);
        window.addEventListener('resize', handleReposition);
        window.addEventListener('scroll', handleReposition, true);

        return () => {
            document.removeEventListener('pointerdown', handlePointerDown);
            window.removeEventListener('resize', handleReposition);
            window.removeEventListener('scroll', handleReposition, true);
        };
    });
</script>

<button
    bind:this={triggerEl}
    type="button"
    class="btn btn-ghost btn-sm gap-1"
    aria-label="Select book and chapter"
    aria-expanded={open}
    onclick={togglePicker}
>
    <span class="max-w-40 truncate">{locationLabel}</span>
    <ChevronDown size={14} aria-hidden="true" />
</button>

{#if open}
    <div
        use:portal
        bind:this={panelEl}
        class="bg-base-100 rounded-box fixed z-[1000] w-[min(36rem,calc(100vw-1.5rem))] border border-base-300 p-3 shadow-lg"
        style={panelStyle}
        role="dialog"
        aria-label="Book and chapter picker"
    >
        {#if step === 'chapter'}
            <button
                type="button"
                class="btn btn-ghost btn-xs mb-2 gap-1 px-1"
                onclick={goBack}
            >
                <ArrowLeft size={14} aria-hidden="true" />
                Back
            </button>
        {/if}

        {#if step === 'book'}
            {#if !booksReady || bible.booksLoading}
                <div class="flex items-center justify-center py-8">
                    <span
                        class="loading loading-spinner loading-md text-primary"
                    ></span>
                </div>
            {:else}
                <div class="grid grid-cols-2 gap-3">
                    <section class="min-w-0">
                        <p class="menu-title px-0">Old Testament</p>
                        <ul
                            class="menu menu-sm rounded-box bg-base-200 max-h-64 overflow-y-auto p-1"
                        >
                            {#each otBooks as book (book.id)}
                                <li>
                                    <button
                                        type="button"
                                        disabled={!book.available}
                                        class:menu-active={study.bookId ===
                                            book.id}
                                        class:opacity-40={!book.available}
                                        onclick={() => selectBook(book.id)}
                                    >
                                        {book.name}
                                    </button>
                                </li>
                            {/each}
                        </ul>
                    </section>
                    <section class="min-w-0">
                        <p class="menu-title px-0">New Testament</p>
                        <ul
                            class="menu menu-sm rounded-box bg-base-200 max-h-64 overflow-y-auto p-1"
                        >
                            {#each ntBooks as book (book.id)}
                                <li>
                                    <button
                                        type="button"
                                        disabled={!book.available}
                                        class:menu-active={study.bookId ===
                                            book.id}
                                        class:opacity-40={!book.available}
                                        onclick={() => selectBook(book.id)}
                                    >
                                        {book.name}
                                    </button>
                                </li>
                            {/each}
                        </ul>
                    </section>
                </div>
            {/if}
        {:else}
            <p class="menu-title px-0">{pendingBook?.name ?? 'Chapter'}</p>
            <div class="grid max-h-64 grid-cols-6 gap-1 overflow-y-auto">
                {#each chapterOptions as chapterNumber (chapterNumber)}
                    <button
                        type="button"
                        class="btn btn-xs"
                        class:btn-primary={pendingBookId === study.bookId &&
                            study.chapter === chapterNumber}
                        onclick={() => selectChapter(chapterNumber)}
                    >
                        {chapterNumber}
                    </button>
                {/each}
            </div>
        {/if}
    </div>
{/if}
