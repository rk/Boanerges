<script lang="ts">
    import { page } from '@inertiajs/svelte';
    import type { Snippet } from 'svelte';

    import SaveStatus from '@/components/scribe/SaveStatus.svelte';
    import ReadabilitySettings from '@/components/settings/ReadabilitySettings.svelte';
    import StudyToolbar from '@/components/toolbar/StudyToolbar.svelte';
    import { bible, loadBooks, loadTranslations } from '@/lib/bible.svelte.ts';
    import { hydrateReadability } from '@/lib/readability.svelte.ts';
    import { closeSettings, hydrateStudy, study } from '@/lib/study.svelte.ts';
    import type { ReadabilitySettings as ReadabilitySettingsType } from '@/lib/types/readability';
    import type { StudySettings as StudySettingsType } from '@/lib/types/study';

    let { children }: { children: Snippet } = $props();

    $effect(() => {
        hydrateReadability(page.props.readability as ReadabilitySettingsType);
        hydrateStudy(page.props.study as StudySettingsType);
    });

    $effect(() => {
        const translationId = study.translationId;

        loadTranslations().then(() => loadBooks(translationId));
    });

    $effect(() => {
        if (bible.translationManagerOpen) {
            return;
        }

        if (! bible.translationsLoaded) {
            return;
        }

        void loadBooks(study.translationId);
    });
</script>

<div class="flex h-dvh min-h-0 flex-col">
    <StudyToolbar />

    <main class="min-h-0 flex-1 overflow-hidden">
        {@render children()}
    </main>
</div>

{#if study.settingsOpen}
    <ReadabilitySettings onclose={closeSettings} />
{/if}

<SaveStatus />
