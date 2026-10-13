import { events, TranslationService, type CatalogItem, type Payload } from "../mygo";

export const translations = $state({
  managerOpen: false,
  catalog: [] as CatalogItem[],
  loading: false,
  managerError: null as string | null,
  /** Modules with an in-flight install click (until API returns). */
  pendingModules: new Set<string>(),
});

export function isTranslationInstalling(entry: CatalogItem): boolean {
  if (translations.pendingModules.has(entry.module)) {
    return true;
  }
  const status = entry.install_status;
  return !!status && status !== "ready" && status !== "failed";
}

function patchCatalogProgress(payload: Payload): void {
  const abbrev = payload.abbrev.toLowerCase();
  translations.catalog = translations.catalog.map((entry) => {
    if (entry.abbrev.toLowerCase() !== abbrev && entry.id.toLowerCase() !== abbrev) {
      return entry;
    }
    const installStatus = payload.install_status ?? payload.step;
    const installed = installStatus === "ready" || payload.step === "ready";
    return {
      ...entry,
      installed,
      install_status: installStatus,
    };
  });
}

let progressUnsub: (() => void) | null = null;

export function bindTranslationInstallProgress(): () => void {
  if (progressUnsub) {
    return progressUnsub;
  }
  progressUnsub = events.translationInstallProgress.on((payload) => {
    patchCatalogProgress(payload);
    for (const entry of translations.catalog) {
      if (entry.abbrev.toLowerCase() === payload.abbrev.toLowerCase()) {
        translations.pendingModules.delete(entry.module);
      }
    }
    if (payload.step === "ready" || payload.step === "failed" || payload.install_status === "ready") {
      void refreshCatalog();
    }
  });
  return progressUnsub;
}

export async function openTranslationManager(): Promise<void> {
  bindTranslationInstallProgress();
  translations.managerOpen = true;
  translations.managerError = null;
  await refreshCatalog();
}

export function closeTranslationManager(): void {
  translations.managerOpen = false;
  translations.managerError = null;
}

export async function refreshCatalog(): Promise<void> {
  translations.loading = true;
  try {
    const resp = await TranslationService.catalog();
    translations.catalog = resp.translations;
  } finally {
    translations.loading = false;
  }
}

export async function installTranslation(module: string): Promise<void> {
  const entry = translations.catalog.find((e) => e.module === module);
  if (translations.pendingModules.has(module) || (entry && isTranslationInstalling(entry))) {
    return;
  }
  translations.pendingModules.add(module);
  translations.managerError = null;
  try {
    await TranslationService.install(module);
    await refreshCatalog();
  } catch (err) {
    translations.pendingModules.delete(module);
    translations.managerError = err instanceof Error ? err.message : String(err);
  }
}

export async function uninstallTranslation(module: string): Promise<void> {
  translations.managerError = null;
  try {
    await TranslationService.uninstall(module);
    await refreshCatalog();
  } catch (err) {
    translations.managerError = err instanceof Error ? err.message : String(err);
  }
}
