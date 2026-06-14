<script lang="ts">
    import { study, setView } from '@/lib/study.svelte.ts';
    import type { ViewMode } from '@/lib/types/bible';

    const views: { id: ViewMode; label: string }[] = [
        { id: 'bible', label: 'Bible' },
        { id: 'comparison', label: 'Comparison' },
        { id: 'scribe', label: 'Scribe' },
    ];
</script>

<div class="join grid w-full grid-cols-3">
    {#each views as view (view.id)}
        <button
            type="button"
            class="btn btn-sm join-item"
            class:btn-primary={study.activeView === view.id}
            onclick={() => setView(view.id)}
        >
            {view.label}
        </button>
    {/each}
</div>
