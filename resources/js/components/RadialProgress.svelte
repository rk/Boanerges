<script lang="ts">
    import { onMount } from 'svelte';

    const Native = window.Native;

    let { abbrev } : { abbrev: string } = $props();
    let progress = $state(0);

    onMount(() => {
        Native.on('App\\Events\\TranslationInstallProgress', payload => {
            if (abbrev.toLowerCase() === (payload.abbrev ?? '').toLowerCase()) {
                progress = payload.percent ?? 0;
            }
        });
    });
</script>

<div class="radial-progress text-primary"
     style="--value:{progress}; --size: 30px; thickness: 2px;"
     aria-valuenow={progress}
     role="progressbar">{progress}%</div>
