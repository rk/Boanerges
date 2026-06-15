<script lang="ts">
    import AppHead from '@/components/AppHead.svelte';
    import AppLayout from '@/layouts/AppLayout.svelte';
    import { study } from '@/lib/study.svelte.ts';
    import BibleView from '@/views/BibleView.svelte';
    import ComparisonView from '@/views/ComparisonView.svelte';
    import ScribeView from '@/views/ScribeView.svelte';
</script>

<AppHead title="Study" />

<AppLayout>
    {#if study.activeView === 'bible'}
        <BibleView />
    {:else if study.activeView === 'comparison'}
        <ComparisonView />
    {:else}
        <ScribeView />
    {/if}
</AppLayout>
