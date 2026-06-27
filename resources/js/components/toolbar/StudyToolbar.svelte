<script lang="ts">
    import { Settings, BookPlus } from '@lucide/svelte/icons';

    import ScrollSyncToggle from '@/components/reader/ScrollSyncToggle.svelte';
    import TranslationManagerModal from '@/components/sidebar/TranslationManagerModal.svelte';
    import BookChapterPicker from '@/components/toolbar/BookChapterPicker.svelte';
    import { openTranslationManager } from '@/lib/bible.svelte.ts';
    import { bibleColumnCount } from '@/lib/studyLayout';
    import {
        openSettings,
        setColumnCount,
        study,
    } from '@/lib/study.svelte.ts';

    const showScrollSync = $derived(bibleColumnCount(study) >= 2);
</script>

<header class="navbar bg-base-100 border-base-300 shrink-0 gap-2 border-b px-3">
    <span class="text-lg font-semibold shrink-0">Boanerges</span>

    <div class="divider divider-horizontal mx-0 hidden h-full sm:flex"></div>

    <BookChapterPicker />

    <button
        type="button"
        class="btn btn-ghost btn-square btn-sm"
        aria-label="Manage translations"
        onclick={openTranslationManager}
    >
        <BookPlus size={16} aria-hidden="true" />
    </button>

    <div class="divider divider-horizontal mx-0 hidden h-full md:flex"></div>

    <div class="join hidden md:flex">
        {#each [1, 2, 3] as count (count)}
            <button
                type="button"
                class="btn btn-sm join-item"
                class:btn-active={study.columnCount === count}
                aria-pressed={study.columnCount === count}
                onclick={() => setColumnCount(count as 1 | 2 | 3)}
            >
                {count}
            </button>
        {/each}
    </div>

    <div class="flex-1"></div>

    {#if showScrollSync}
        <ScrollSyncToggle />
    {/if}

    <button
        type="button"
        class="btn btn-ghost btn-square btn-sm"
        aria-label="Edit settings"
        onclick={openSettings}
    >
        <Settings size={18} aria-hidden="true" />
    </button>
</header>

<TranslationManagerModal />
