<script lang="ts">
  import {
    closeTranslationManager,
    installTranslation,
    translations,
    uninstallTranslation,
  } from "../../lib/translations.svelte.ts";

  let dialog = $state<HTMLDialogElement | null>(null);
  let query = $state("");

  const filtered = $derived(
    translations.catalog.filter((entry) => {
      const needle = query.trim().toLowerCase();
      if (!needle) return true;
      return entry.name.toLowerCase().includes(needle) || entry.abbrev.toLowerCase().includes(needle);
    }),
  );

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
        <li class="flex items-center justify-between gap-3 py-3">
          <div class="min-w-0">
            <p class="font-medium truncate">{entry.name}</p>
            <p class="text-xs opacity-70">{entry.abbrev}{entry.bundled ? " · bundled" : ""}</p>
          </div>
          {#if entry.installed && !entry.bundled}
            <button type="button" class="btn btn-sm btn-outline" onclick={() => uninstallTranslation(entry.module)}>Remove</button>
          {:else if !entry.installed}
            <button type="button" class="btn btn-sm btn-primary" onclick={() => installTranslation(entry.module)}>Install</button>
          {:else}
            <span class="badge badge-ghost">Installed</span>
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
