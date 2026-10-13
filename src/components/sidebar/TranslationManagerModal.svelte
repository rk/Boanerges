<script lang="ts">
  import CircleHelp from "@lucide/svelte/icons/circle-help";
  import type { CatalogItem } from "../../mygo";
  import {
    catalogProgress,
    closeTranslationManager,
    installStepLabel,
    installTranslation,
    isTranslationInstalling,
    translations,
    uninstallTranslation,
  } from "../../lib/translations.svelte.ts";

  let dialog = $state<HTMLDialogElement | null>(null);
  let infoDialog = $state<HTMLDialogElement | null>(null);
  let query = $state("");
  let infoEntry = $state<CatalogItem | null>(null);

  const filtered = $derived(
    translations.catalog.filter((entry) => {
      const needle = query.trim().toLowerCase();
      if (!needle) return true;
      return entry.name.toLowerCase().includes(needle) || entry.abbrev.toLowerCase().includes(needle);
    }),
  );

  function openInfo(entry: CatalogItem, event: MouseEvent) {
    event.stopPropagation();
    infoEntry = entry;
    infoDialog?.showModal();
  }

  function closeInfo() {
    infoEntry = null;
    infoDialog?.close();
  }

  function aboutIsUrl(text: string): boolean {
    return /^https?:\/\//i.test(text.trim());
  }

  $effect(() => {
    if (translations.managerOpen) {
      dialog?.showModal();
    } else {
      dialog?.close();
    }
  });
</script>

<dialog bind:this={dialog} class="modal" data-testid="translation-manager" onclose={closeTranslationManager}>
  <div class="modal-box flex max-h-[85vh] max-w-lg flex-col">
    <h3 class="text-lg font-bold">Manage Translations</h3>
    <input
      class="input input-bordered mt-4"
      type="search"
      placeholder="Search translations…"
      bind:value={query}
      data-testid="translation-search"
    />
    {#if translations.managerError}
      <div class="alert alert-error mt-4 text-sm" data-testid="translation-manager-error">
        <span>{translations.managerError}</span>
      </div>
    {/if}
    <ul class="mt-4 min-h-0 flex-1 overflow-y-auto divide-y divide-base-300">
      {#each filtered as entry (entry.module)}
        {@const installing = isTranslationInstalling(entry)}
        {@const progress = catalogProgress(entry)}
        <li class="flex items-start justify-between gap-3 py-3">
          <div class="min-w-0 flex-1">
            <div class="flex items-center gap-1">
              <p class="font-medium truncate">{entry.name}</p>
              <button
                type="button"
                class="btn btn-ghost btn-xs btn-square shrink-0"
                aria-label="Translation info"
                data-testid="translation-info"
                onclick={(e) => openInfo(entry, e)}
              >
                <CircleHelp class="size-4 opacity-70" />
              </button>
            </div>
            <p class="text-xs opacity-70">
              {entry.abbrev}{entry.bundled ? " · bundled" : ""}
            </p>
            {#if installing && progress}
              <p class="text-xs mt-1" data-testid="translation-install-step">
                {installStepLabel(progress.step)}
              </p>
            {/if}
          </div>
          {#if entry.installed && !entry.bundled}
            <button type="button" class="btn btn-sm btn-outline shrink-0" onclick={() => uninstallTranslation(entry.module)}>
              Remove
            </button>
          {:else if installing}
            <button type="button" class="btn btn-sm btn-primary shrink-0" disabled data-testid="translation-installing">
              <span class="loading loading-spinner loading-xs"></span>
              {progress ? `${progress.percent}%` : "…"}
            </button>
          {:else if !entry.installed}
            <button type="button" class="btn btn-sm btn-primary shrink-0" onclick={() => installTranslation(entry.module)}>
              Install
            </button>
          {:else}
            <span class="badge badge-ghost shrink-0">Installed</span>
          {/if}
        </li>
      {/each}
    </ul>
    <div class="modal-action">
      <button type="button" class="btn" onclick={closeTranslationManager}>Close</button>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button type="submit">close</button></form>
</dialog>

<dialog bind:this={infoDialog} class="modal" data-testid="translation-info-dialog" onclose={closeInfo}>
  <div class="modal-box max-w-md">
    {#if infoEntry}
      <h3 class="text-lg font-bold">{infoEntry.name}</h3>
      <p class="text-sm opacity-70 mt-1">{infoEntry.abbrev}{infoEntry.bundled ? " · bundled" : ""}</p>
      <div class="mt-4 text-sm">
        {#if infoEntry.about?.trim()}
          {#if aboutIsUrl(infoEntry.about)}
            <a class="link link-primary" href={infoEntry.about} target="_blank" rel="noopener noreferrer">
              Module information
            </a>
          {:else}
            <p class="whitespace-pre-wrap">{infoEntry.about}</p>
          {/if}
        {:else}
          <p class="opacity-60">No description in catalog.</p>
        {/if}
      </div>
    {/if}
    <div class="modal-action">
      <form method="dialog">
        <button type="submit" class="btn">Close</button>
      </form>
    </div>
  </div>
  <form method="dialog" class="modal-backdrop"><button type="submit">close</button></form>
</dialog>
