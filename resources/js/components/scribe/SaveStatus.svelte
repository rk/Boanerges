<script lang="ts">
    import { scribe } from '@/lib/scribe.svelte.ts';
</script>

{#if scribe.saveStatus !== 'idle'}
    <div class="toast toast-end toast-bottom z-50">
        <div class="alert alert-success alert-soft">
            {#if scribe.saveStatus === 'saving'}
                Saving…
            {:else}
                Saved
            {/if}
        </div>
    </div>
{/if}
