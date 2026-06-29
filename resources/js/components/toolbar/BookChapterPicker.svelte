<script lang="ts">
    import ArrowLeft from '@lucide/svelte/icons/arrow-left';
    import ChevronDown from '@lucide/svelte/icons/chevron-down';

    import { bible } from '@/lib/bible.svelte.ts';
    import { study, setBook, setChapter } from '@/lib/study.svelte.ts';
    import type { Testament } from '@/lib/types/bible';

    type Step = 'testament' | 'book' | 'chapter';

    let triggerEl = $state<HTMLButtonElement | null>(null);
    let panelEl = $state<HTMLDivElement | null>(null);
    let open = $state(false);
    let panelStyle = $state('');
    let step = $state<Step>('testament');
    let selectedTestament = $state<Testament | null>(null);
    let pendingBookId = $state<string | null>(null);

    const currentBook = $derived(
        bible.books.find((book) => book.id === study.bookId),
    );
    const locationLabel = $derived(
        currentBook ? `${currentBook.name} ${study.chapter}` : 'Book & chapter',
    );

    const booksInTestament = $derived(
        bible.books.filter((book) => book.testament === selectedTestament),
    );

    const pendingBook = $derived(
        bible.books.find((book) => book.id === pendingBookId) ?? currentBook,
    );

    const chapterOptions = $derived(
        Array.from(
            { length: pendingBook?.chapters ?? 1 },
            (_, index) => index + 1,
        ),
    );

    function resetPicker(): void {
        step = 'testament';
        selectedTestament = null;
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

    function selectTestament(testament: Testament): void {
        selectedTestament = testament;
        step = 'book';
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
        if (step === 'chapter') {
            step = 'book';
            pendingBookId = null;

            return;
        }

        if (step === 'book') {
            step = 'testament';
            selectedTestament = null;
        }
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
        class="bg-base-100 rounded-box fixed z-[1000] w-72 border border-base-300 p-3 shadow-lg"
        style={panelStyle}
        role="dialog"
        aria-label="Book and chapter picker"
    >
        {#if step !== 'testament'}
            <button
                type="button"
                class="btn btn-ghost btn-xs mb-2 gap-1 px-1"
                onclick={goBack}
            >
                <ArrowLeft size={14} aria-hidden="true" />
                Back
            </button>
        {/if}

        {#if step === 'testament'}
            <p class="menu-title px-0">Testament</p>
            <div class="join grid w-full grid-cols-2">
                <button
                    type="button"
                    class="btn btn-sm join-item"
                    class:btn-primary={currentBook?.testament === 'ot'}
                    onclick={() => selectTestament('ot')}
                >
                    Old Testament
                </button>
                <button
                    type="button"
                    class="btn btn-sm join-item"
                    class:btn-primary={currentBook?.testament === 'nt'}
                    onclick={() => selectTestament('nt')}
                >
                    New Testament
                </button>
            </div>
        {:else if step === 'book'}
            <p class="menu-title px-0">
                {selectedTestament === 'ot' ? 'Old Testament' : 'New Testament'}
            </p>
            {#if bible.booksLoading}
                <p class="text-base-content/60 px-1 py-2 text-sm">Loading…</p>
            {:else}
                <ul
                    class="menu menu-sm rounded-box bg-base-200 max-h-64 overflow-y-auto p-1"
                >
                    {#each booksInTestament as book (book.id)}
                        <li>
                            <button
                                type="button"
                                class:menu-active={study.bookId === book.id}
                                onclick={() => selectBook(book.id)}
                            >
                                {book.name}
                            </button>
                        </li>
                    {/each}
                </ul>
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
