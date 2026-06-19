<script type="ts">
    import { Menu, Search, NotebookPen, Printer, BookMarked, MessageSquareText } from '@lucide/svelte';
    import { openSearch } from '@/lib/search.svelte.ts';
    import { openCrossReferences } from '@/lib/crossrefs.svelte.ts';
    import { study } from '@/lib/study.svelte.ts';
</script>

<div class="fab">
    <div tabindex="0" role="button" class="btn btn-lg btn-circle btn-info">
        <Menu size="18" color="currentColor" />
    </div>

    <div class="fab-close">
        <span class="btn btn-circle btn-lg btn-error">✕</span>
    </div>

    <div class="tooltip tooltip-left" data-tip="Search">
        <button type="button" class="btn btn-lg btn-circle" onclick={openSearch}><Search /></button>
    </div>
    <div class="tooltip tooltip-left" data-tip="Notes">
        <button type="button" class="btn btn-lg btn-circle"><NotebookPen /></button>
    </div>
    <div class="tooltip tooltip-left" data-tip="Cross-References">
        <button
            type="button"
            class="btn btn-lg btn-circle"
            onclick={() => openCrossReferences(study.chapter > 0 ? 1 : 1)}
        >
            <BookMarked />
        </button>
    </div>
    <div class="tooltip tooltip-left" data-tip="Print">
        <button type="button" class="btn btn-lg btn-circle"><Printer /></button>
    </div>
    <div class="tooltip tooltip-left" data-tip="Feedback">
        <button type="button" class="btn btn-lg btn-circle"><MessageSquareText /></button>
    </div>
</div>
