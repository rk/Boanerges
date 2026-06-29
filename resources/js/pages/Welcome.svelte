<script lang="ts">
    import { router } from '@inertiajs/svelte';
    import { onMount } from 'svelte';
    import { watchInstallProgress } from '@/lib/nativeBroadcast.ts';

    let progress = $state({
        step: 'starting',
        percent: 0,
        install_status: 'pending',
    });

    onMount(() => {
        return watchInstallProgress(
            'App\\Events\\TranslationInstallProgress',
            '/bible/translations/ASV/install-status',
            (payload) => {
                progress = {
                    step: payload.step ?? progress.step,
                    percent: payload.percent ?? progress.percent,
                    install_status:
                        payload.install_status ??
                        progress.install_status ??
                        'pending',
                };

                if (
                    payload.step === 'ready' ||
                    payload.install_status === 'ready'
                ) {
                    router.visit('/');
                }
            },
        );
    });
</script>

<div class="flex min-h-dvh flex-col items-center justify-center gap-6 p-8">
    <div class="text-center">
        <h1 class="text-3xl font-bold">Welcome to Boanerges</h1>
        <p class="text-base-content/70 mt-2 max-w-md">
            Setting up the bundled American Standard Version translation. This
            runs once on first launch.
        </p>
    </div>

    <div class="w-full max-w-md">
        <progress
            class="progress progress-primary w-full"
            value={progress.percent}
            max="100"
        ></progress>
        <p class="text-base-content/60 mt-2 text-center text-sm capitalize">
            {progress.step.replaceAll('_', ' ')}
        </p>
    </div>
</div>
