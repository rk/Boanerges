import { events, TranslationService, type CatalogItem, type Payload } from "../mygo";

export const translations = $state({
  managerOpen: false,
  catalog: [] as CatalogItem[],
  loading: false,
  managerError: null as string | null,
  /** Modules with an in-flight install click (until API returns). */
  pendingModules: new Set<string>(),
  /** Live progress from install events (overrides catalog until refresh). */
  liveProgress: new Map<string, { step: string; percent: number; status: string }>(),
});

export function installStepLabel(step: string | undefined | null): string {
  if (!step) return "Starting…";
  const labels: Record<string, string> = {
    pending: "Queued",
    queued: "Queued",
    downloading: "Downloading",
    downloaded: "Downloaded",
    source_ready: "Source ready",
    creating_schema: "Creating database",
    importing: "Importing verses",
    verifying: "Verifying",
    indexing_coverage: "Indexing books",
    backfill_plain_text: "Preparing search text",
    indexing: "Building search index",
    indexed: "Finishing",
    ready: "Ready",
    failed: "Failed",
  };
  return labels[step] ?? step.replaceAll("_", " ");
}

export function catalogProgress(entry: CatalogItem): { step: string; percent: number; status: string } | null {
  const live = translations.liveProgress.get(entry.module);
  if (live) return live;
  if (!entry.install_status || entry.install_status === "ready" || entry.install_status === "failed") {
    return null;
  }
  return {
    step: entry.install_step ?? entry.install_status,
    percent: entry.install_percent ?? 0,
    status: entry.install_status,
  };
}

export function isTranslationInstalling(entry: CatalogItem): boolean {
  if (translations.pendingModules.has(entry.module)) {
    return true;
  }
  const status = entry.install_status;
  return !!status && status !== "ready" && status !== "failed";
}

function patchCatalogProgress(payload: Payload): void {
  const abbrev = payload.abbrev.toLowerCase();
  const step = payload.step;
  const ready = step === "ready" || payload.install_status === "ready";
  const installStatus = ready ? "ready" : (payload.install_status ?? step);
  const percent = ready ? 100 : (payload.percent ?? 0);

  translations.catalog = translations.catalog.map((entry) => {
    if (entry.abbrev.toLowerCase() !== abbrev && entry.id.toLowerCase() !== abbrev) {
      return entry;
    }
    if (ready) {
      translations.liveProgress.delete(entry.module);
    } else {
      translations.liveProgress.set(entry.module, { step, percent, status: installStatus });
    }
    return {
      ...entry,
      installed: ready || entry.installed,
      install_status: ready ? "ready" : installStatus,
      install_step: ready ? "ready" : step,
      install_percent: percent,
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
      for (const entry of translations.catalog) {
        if (entry.abbrev.toLowerCase() === payload.abbrev.toLowerCase()) {
          translations.liveProgress.delete(entry.module);
        }
      }
      void refreshCatalog();
    }
  });
  return progressUnsub;
}

export async function openTranslationManager(): Promise<void> {
  bindTranslationInstallProgress();
  translations.managerOpen = true;
  translations.managerError = null;
  translations.liveProgress.clear();
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
  translations.liveProgress.set(module, { step: "queued", percent: 5, status: "pending" });
  translations.managerError = null;
  try {
    await TranslationService.install(module);
    await refreshCatalog();
  } catch (err) {
    translations.pendingModules.delete(module);
    translations.liveProgress.delete(module);
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
