<script lang="ts">
    import BookOpen from '@lucide/svelte/icons/book-open';
    import Columns2 from '@lucide/svelte/icons/columns-2';
    import PenLine from '@lucide/svelte/icons/pen-line';

    import { study, setView } from '@/lib/study.svelte.ts';
    import type { ViewMode } from '@/lib/types/bible';

    const views: { id: ViewMode; label: string; icon: typeof BookOpen }[] = [
        { id: 'bible', label: 'Bible', icon: BookOpen },
        { id: 'comparison', label: 'Comparison', icon: Columns2 },
        { id: 'scribe', label: 'Scribe', icon: PenLine },
    ];
</script>

<div class="join grid w-full grid-cols-3">
    {#each views as view (view.id)}
        <button
            type="button"
            class="btn btn-sm join-item gap-1"
            class:btn-primary={study.activeView === view.id}
            onclick={() => setView(view.id)}
        >
            <view.icon size={14} aria-hidden="true" />
            <span>{view.label}</span>
        </button>
    {/each}
</div>
