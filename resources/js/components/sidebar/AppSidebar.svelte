<script lang="ts">
    import BookSelector from '@/components/sidebar/BookSelector.svelte';
    import ChapterSelect from '@/components/sidebar/ChapterSelect.svelte';
    import TranslationSelect from '@/components/sidebar/TranslationSelect.svelte';
    import ViewSelector from '@/components/sidebar/ViewSelector.svelte';
    import {
        study,
        openSettings,
        setTranslation,
        setTranslationB,
    } from '@/lib/study.svelte.ts';

    let { drawerId }: { drawerId: string } = $props();
</script>

<div
    class="is-drawer-close:w-14 is-drawer-open:w-72 bg-base-200 text-base-content flex min-h-full flex-col"
>
    <div class="border-base-300 flex items-center justify-between border-b p-2">
        <span class="is-drawer-close:hidden px-2 text-sm font-semibold">Navigation</span>
        <label
            for={drawerId}
            class="btn btn-ghost btn-circle btn-sm drawer-button is-drawer-open:rotate-y-180 is-drawer-close:tooltip is-drawer-close:tooltip-right"
            aria-label="Toggle sidebar"
            data-tip="Sidebar"
        >
            ↔
        </label>
    </div>

    <div class="is-drawer-close:hidden flex-1 overflow-y-auto p-3">
        <section class="mb-4">
            <p class="menu-title px-0">View</p>
            <ViewSelector />
        </section>

        <section class="mb-4">
            <p class="menu-title px-0">Translation</p>
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

        <ul class="menu rounded-box bg-base-100">
            <li>
                <button type="button" onclick={openSettings}>Readability settings</button>
            </li>
        </ul>
    </div>

    <div class="is-drawer-open:hidden border-base-300 border-t p-2">
        <ul class="menu menu-vertical w-full">
            <li>
                <button
                    type="button"
                    class="is-drawer-close:tooltip is-drawer-close:tooltip-right"
                    data-tip="Settings"
                    onclick={openSettings}
                >
                    ⚙
                </button>
            </li>
        </ul>
    </div>
</div>
