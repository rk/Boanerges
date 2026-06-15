<script lang="ts">
    import { page } from '@inertiajs/svelte';
    import PanelLeftClose from '@lucide/svelte/icons/panel-left-close';
    import type { Snippet } from 'svelte';

    import SaveStatus from '@/components/scribe/SaveStatus.svelte';
    import ReadabilitySettings from '@/components/settings/ReadabilitySettings.svelte';
    import AppSidebar from '@/components/sidebar/AppSidebar.svelte';
    import { loadBooks, loadTranslations } from '@/lib/bible.svelte.ts';
    import { hydrateReadability } from '@/lib/readability.svelte.ts';
    import { closeSettings, hydrateStudy, study } from '@/lib/study.svelte.ts';
    import type { ReadabilitySettings as ReadabilitySettingsType } from '@/lib/types/readability';
    import type { StudySettings as StudySettingsType } from '@/lib/types/study';

    let { children }: { children: Snippet } = $props();

    const drawerId = 'app-drawer';

    $effect(() => {
        hydrateReadability(page.props.readability as ReadabilitySettingsType);
        hydrateStudy(page.props.study as StudySettingsType);
    });

    $effect(() => {
        const translationId = study.translationId;

        loadTranslations().then(() => loadBooks(translationId));
    });
</script>

<div class="drawer lg:drawer-open h-dvh">
    <input id={drawerId} type="checkbox" class="drawer-toggle" />

    <div class="drawer-content flex min-h-0 flex-col">
        <header class="navbar bg-base-100 border-base-300 shrink-0 border-b px-4">
            <label
                for={drawerId}
                class="btn btn-ghost btn-square drawer-button"
                aria-label="Toggle sidebar"
            >
                <PanelLeftClose size={20} aria-hidden="true" />
            </label>
            <span class="text-lg font-semibold">Boanerges</span>
        </header>

        <main class="min-h-0 flex-1 overflow-hidden">
            {@render children()}
        </main>
    </div>

    <div class="drawer-side is-drawer-close:overflow-visible z-20">
        <label for={drawerId} aria-label="Close sidebar" class="drawer-overlay lg:hidden"></label>
        <AppSidebar {drawerId} />
    </div>
</div>

{#if study.settingsOpen}
    <ReadabilitySettings onclose={closeSettings} />
{/if}

<SaveStatus />
