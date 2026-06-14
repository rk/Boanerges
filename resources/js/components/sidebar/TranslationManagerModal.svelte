<script lang="ts">
    import {
        bible,
        closeTranslationManager,
        installTranslation,
        uninstallTranslation,
    } from '@/lib/bible.svelte.ts';
    import { syncStudyTranslationSelection } from '@/lib/study.svelte.ts';
    import type { CatalogTranslation } from '@/lib/types/bible';

    let dialog: HTMLDialogElement | undefined = $state();
    let query = $state('');

    const filtered = $derived(
        bible.catalog.filter((entry) => {
            const needle = query.trim().toLowerCase();

            if (needle === '') {
                return true;
            }

            return (
                entry.name.toLowerCase().includes(needle)
                || entry.abbrev.toLowerCase().includes(needle)
            );
        }),
    );

    $effect(() => {
        if (bible.translationManagerOpen) {
            dialog?.showModal();
        } else {
            dialog?.close();
        }
    });

    function handleClose(): void {
        closeTranslationManager();
    }

    async function handleInstall(entry: CatalogTranslation): Promise<void> {
        await installTranslation(entry.module);
        syncStudyTranslationSelection();
    }

    async function handleUninstall(entry: CatalogTranslation): Promise<void> {
        await uninstallTranslation(entry.module);
        syncStudyTranslationSelection();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box flex max-h-[85vh] max-w-lg flex-col">
        <h3 class="text-lg font-bold">Manage Translations</h3>
        <p class="text-base-content/70 mt-1 text-sm">
            English translations available for download. Bundled translations ship with the app.
        </p>

        <label class="input input-bordered mt-4 flex items-center gap-2">
            <input
                type="search"
                class="grow"
                placeholder="Search translations…"
                bind:value={query}
            />
        </label>

        {#if bible.managerError}
            <div class="alert alert-error mt-4 text-sm">
                <span>{bible.managerError}</span>
            </div>
        {/if}

        <div class="mt-4 min-h-0 flex-1 overflow-y-auto">
            {#if bible.catalogLoading}
                <div class="flex justify-center py-8">
                    <span class="loading loading-spinner loading-md text-primary"></span>
                </div>
            {:else}
                <ul class="divide-base-300 divide-y">
                    {#each filtered as entry (entry.module)}
                        <li class="flex items-start justify-between gap-3 py-3">
                            <div class="min-w-0">
                                <p class="truncate font-medium">{entry.name}</p>
                                <div class="mt-1 flex flex-wrap gap-1">
                                    <span class="badge badge-outline badge-sm">{entry.abbrev}</span>
                                    {#if entry.bundled}
                                        <span class="badge badge-neutral badge-sm">Bundled</span>
                                    {/if}
                                    {#if entry.installed && ! entry.bundled}
                                        <span class="badge badge-success badge-sm">Installed</span>
                                    {/if}
                                </div>
                            </div>

                            <div class="shrink-0">
                                {#if entry.bundled}
                                    <span class="text-base-content/50 text-xs">Included</span>
                                {:else if entry.installed}
                                    <button
                                        type="button"
                                        class="btn btn-ghost btn-xs"
                                        disabled={bible.uninstallingModule === entry.module}
                                        onclick={() => handleUninstall(entry)}
                                    >
                                        {#if bible.uninstallingModule === entry.module}
                                            <span class="loading loading-spinner loading-xs"></span>
                                        {:else}
                                            Remove
                                        {/if}
                                    </button>
                                {:else}
                                    <button
                                        type="button"
                                        class="btn btn-primary btn-xs"
                                        disabled={bible.installingModule === entry.module}
                                        onclick={() => handleInstall(entry)}
                                    >
                                        {#if bible.installingModule === entry.module}
                                            <span class="loading loading-spinner loading-xs"></span>
                                        {:else}
                                            Install
                                        {/if}
                                    </button>
                                {/if}
                            </div>
                        </li>
                    {/each}
                </ul>
            {/if}
        </div>

        <div class="modal-action">
            <form method="dialog">
                <button type="submit" class="btn">Close</button>
            </form>
        </div>
    </div>

    <form method="dialog" class="modal-backdrop">
        <button aria-label="Close translation manager">close</button>
    </form>
</dialog>
