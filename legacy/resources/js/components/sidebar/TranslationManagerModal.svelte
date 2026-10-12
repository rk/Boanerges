<script lang="ts">
    import { Trash2, Plus, Search, Info } from '@lucide/svelte';
    import RadialProgress from '@/components/RadialProgress.svelte';
    import {
        bible,
        closeTranslationManager,
        installTranslation,
        loadTranslations,
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
                entry.name.toLowerCase().includes(needle) ||
                entry.abbrev.toLowerCase().includes(needle)
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

    function isInstalling(entry: CatalogTranslation): boolean {
        if (bible.installingModule === entry.module) {
            return true;
        }

        return Boolean(
            entry.install_status &&
            entry.install_status !== 'ready' &&
            entry.install_status !== 'failed',
        );
    }

    function handleClose(): void {
        closeTranslationManager();
    }

    async function handleInstall(entry: CatalogTranslation): Promise<void> {
        await installTranslation(entry.module);
        await loadTranslations(true);
        syncStudyTranslationSelection();
    }

    async function handleUninstall(entry: CatalogTranslation): Promise<void> {
        await uninstallTranslation(entry.module);
        await loadTranslations(true);
        syncStudyTranslationSelection();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box flex max-h-[85vh] max-w-lg flex-col">
        <h3 class="text-lg font-bold">Manage Translations</h3>
        <p class="text-base-content/70 mt-1 text-sm">
            English translations available for download. Bundled translations
            ship with the app.
        </p>

        <label class="input input-bordered mt-4">
            <Search size="14" />
            <input
                type="search"
                class="grow py-2"
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
            {#if bible.catalogLoading && bible.catalog.length === 0}
                <div class="flex justify-center py-8">
                    <span
                        class="loading loading-spinner loading-md text-primary"
                    ></span>
                </div>
            {:else if filtered.length === 0}
                <p class="text-base-content/60 py-4 text-sm">
                    No translations match your search.
                </p>
            {:else}
                <ul class="divide-base-300 divide-y">
                    {#each filtered as entry (entry.module)}
                        <li
                            class="flex items-center justify-between gap-3 py-3"
                        >
                            <div class="min-w-0">
                                <p
                                    class="font-medium flex items-center flex-nowrap gap-1"
                                >
                                    <span class="truncate">{entry.name}</span>
                                    <a
                                        class="link link-secondary"
                                        href={entry.about}
                                        target="_blank"><Info size="14" /></a
                                    >
                                </p>
                                <div class="mt-1 flex flex-wrap gap-1">
                                    <span class="badge badge-outline badge-sm"
                                        >{entry.abbrev}</span
                                    >
                                    {#if entry.bundled}
                                        <span
                                            class="badge badge-neutral badge-sm"
                                            >Bundled</span
                                        >
                                    {:else if entry.installed}
                                        <span
                                            class="badge badge-success badge-sm"
                                            >Installed</span
                                        >
                                    {:else if isInstalling(entry)}
                                        <span
                                            class="badge badge-warning badge-sm"
                                            >Installing</span
                                        >
                                    {/if}
                                </div>
                            </div>

                            <div class="shrink-0">
                                {#if entry.bundled}
                                    <span class="text-base-content/50 text-xs"
                                        >Included</span
                                    >
                                {:else if entry.installed}
                                    <button
                                        type="button"
                                        class="btn btn-error btn-sm btn-square"
                                        disabled={bible.uninstallingModule ===
                                            entry.module || isInstalling(entry)}
                                        onclick={() => handleUninstall(entry)}
                                    >
                                        {#if bible.uninstallingModule === entry.module}
                                            <span
                                                class="loading loading-spinner loading-xs"
                                            ></span>
                                        {:else}
                                            <div
                                                class="tooltip tooltip-left"
                                                data-tip="Remove"
                                            >
                                                <Trash2 size="16" />
                                            </div>
                                        {/if}
                                    </button>
                                {:else}
                                    <button
                                        type="button"
                                        class="btn btn-primary btn-sm btn-square"
                                        disabled={bible.installingModule !==
                                            null}
                                        onclick={() => handleInstall(entry)}
                                    >
                                        {#if bible.installingModule === entry.module}
                                            <RadialProgress
                                                abbrev={entry.abbrev}
                                            />
                                        {:else}
                                            <span
                                                class="tooltip tooltip-left"
                                                data-tip="Install"
                                            >
                                                <Plus size="16" />
                                            </span>
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
