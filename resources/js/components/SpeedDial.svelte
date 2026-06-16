<script type="ts">
    import { Menu, Search, NotebookPen, Printer, BookMarked, MessageSquareText } from '@lucide/svelte';

</script>

<div class="fab">
    <!-- a focusable div with tabindex is necessary to work on all browsers. role="button" is necessary for accessibility -->
    <div tabindex="0" role="button" class="btn btn-lg btn-circle btn-info">
        <Menu size="18" color="currentColor" />
    </div>

    <!-- close button should not be focusable so it can close the FAB when clicked. It's just a visual placeholder -->
    <div class="fab-close">
        <span class="btn btn-circle btn-lg btn-error">✕</span>
    </div>

    <!-- buttons that show up when FAB is open -->
    <div class="tooltip tooltip-left" data-tip="Search"><button class="btn btn-lg btn-circle"><Search /></button></div>
    <div class="tooltip tooltip-left" data-tip="Notes"><button class="btn btn-lg btn-circle"><NotebookPen /></button></div>
    <div class="tooltip tooltip-left" data-tip="Cross-References"><button class="btn btn-lg btn-circle"><BookMarked /></button></div>
    <div class="tooltip tooltip-left" data-tip="Print"><button class="btn btn-lg btn-circle"><Printer /></button></div>
    <div class="tooltip tooltip-left" data-tip="Feedback"><button class="btn btn-lg btn-circle"><MessageSquareText /></button></div>
</div>
