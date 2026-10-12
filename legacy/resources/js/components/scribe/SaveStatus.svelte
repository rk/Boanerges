<script lang="ts">
    import CircleCheck from '@lucide/svelte/icons/circle-check';
    import LoaderCircle from '@lucide/svelte/icons/loader-circle';

    import { scribe } from '@/lib/scribe.svelte.ts';
</script>

{#if scribe.saveStatus !== 'idle'}
    <div class="toast toast-end toast-bottom z-50">
        <div class="alert alert-success alert-soft">
            {#if scribe.saveStatus === 'saving'}
                <LoaderCircle
                    size={18}
                    class="animate-spin"
                    aria-hidden="true"
                />
                Saving…
            {:else}
                <CircleCheck size={18} aria-hidden="true" />
                Saved
            {/if}
        </div>
    </div>
{/if}
