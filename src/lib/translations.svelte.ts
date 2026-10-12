import { TranslationService, type CatalogItem } from "../mygo";

export const translations = $state({
  managerOpen: false,
  catalog: [] as CatalogItem[],
  loading: false,
  managerError: null as string | null,
});

export async function openTranslationManager(): Promise<void> {
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
  translations.managerError = null;
  try {
    await TranslationService.install(module);
    await refreshCatalog();
  } catch (err) {
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
