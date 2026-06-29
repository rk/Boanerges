<script lang="ts">
    import { onMount } from 'svelte';

    let { abbrev }: { abbrev: string } = $props();
    let progress = $state(0);

    onMount(() => {
        const native = window.Native;

        if (!native) {
            return;
        }

        native.on(
            'App\\Events\\TranslationInstallProgress',
            (payload: unknown) => {
                const data = payload as { abbrev?: string; percent?: number };

                if (
                    abbrev.toLowerCase() === (data.abbrev ?? '').toLowerCase()
                ) {
                    progress = data.percent ?? 0;
                }
            },
        );
    });
</script>

<div
    class="radial-progress text-primary"
    style="--value:{progress}; --size: 30px; thickness: 2px;"
    aria-valuenow={progress}
    role="progressbar"
>
    {progress}%
</div>
