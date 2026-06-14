<script lang="ts">
    import BookOpen from '@lucide/svelte/icons/book-open';
    import Columns2 from '@lucide/svelte/icons/columns-2';
    import PanelLeftClose from '@lucide/svelte/icons/panel-left-close';
    import PenLine from '@lucide/svelte/icons/pen-line';
    import SquarePlus from '@lucide/svelte/icons/square-plus';
    import Type from '@lucide/svelte/icons/type';

    import BookSelector from '@/components/sidebar/BookSelector.svelte';
    import ChapterSelect from '@/components/sidebar/ChapterSelect.svelte';
    import CollapsedBookChapterPicker from '@/components/sidebar/CollapsedBookChapterPicker.svelte';
    import TranslationManagerModal from '@/components/sidebar/TranslationManagerModal.svelte';
    import TranslationSelect from '@/components/sidebar/TranslationSelect.svelte';
    import { openTranslationManager } from '@/lib/bible.svelte.ts';
    import {
        study,
        openSettings,
        setTranslation,
        setTranslationB,
        setView,
    } from '@/lib/study.svelte.ts';
    import type { ViewMode } from '@/lib/types/bible';

    let { drawerId }: { drawerId: string } = $props();

    const views: { id: ViewMode; label: string; icon: typeof BookOpen }[] = [
        { id: 'bible', label: 'Bible', icon: BookOpen },
        { id: 'comparison', label: 'Comparison', icon: Columns2 },
        { id: 'scribe', label: 'Scribe', icon: PenLine },
    ];
</script>

<div
    class="sidebar-shell bg-base-200 text-base-content flex min-h-full w-full flex-col items-start"
>
    <ul class="menu w-full shrink-0">
        {#each views as view (view.id)}
            <li>
                <button
                    type="button"
                    class="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                    class:menu-active={study.activeView === view.id}
                    data-tip={view.label}
                    aria-label={view.label}
                    onclick={() => setView(view.id)}
                >
                    <view.icon size={18} aria-hidden="true" />
                    <span class="is-drawer-close:hidden">{view.label}</span>
                </button>
            </li>
        {/each}

        <CollapsedBookChapterPicker />

        <li>
            <button
                type="button"
                class="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                data-tip="Readability settings"
                aria-label="Readability settings"
                onclick={openSettings}
            >
                <Type size={18} aria-hidden="true" />
                <span class="is-drawer-close:hidden">Readability settings</span>
            </button>
        </li>
    </ul>

    <div class="sidebar-expanded-only w-full flex-1 overflow-y-auto p-3">
        <section class="mb-4">
            <div class="mb-1 flex items-center justify-between gap-2">
                <p class="menu-title px-0">Translation</p>
                <button
                    type="button"
                    class="btn btn-ghost btn-xs btn-square"
                    aria-label="Manage translations"
                    onclick={openTranslationManager}
                >
                    <SquarePlus size={16} aria-hidden="true" />
                </button>
            </div>
            <TranslationSelect label="Primary" value={study.translationId} onchange={setTranslation} />

            {#if study.activeView === 'comparison' || study.activeView === 'scribe'}
                <div class="mt-2">
                    <TranslationSelect label="Secondary" value={study.translationBId} onchange={setTranslationB} />
                </div>
            {/if}
        </section>

        <section class="mb-4">
            <p class="menu-title px-0">Book</p>
            <BookSelector />
        </section>

        <section class="mb-4">
            <ChapterSelect />
        </section>
    </div>

    <div
        class="is-drawer-close:tooltip is-drawer-close:tooltip-right m-2 mt-auto"
        data-tip="Toggle sidebar"
    >
        <label
            for={drawerId}
            class="btn btn-ghost btn-circle drawer-button is-drawer-open:rotate-y-180"
            aria-label="Toggle sidebar"
        >
            <PanelLeftClose size={18} aria-hidden="true" />
        </label>
    </div>
</div>

<TranslationManagerModal />
