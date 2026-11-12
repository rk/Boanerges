export type StudioProjectKind = "dictionary" | "topical";

const keys: Record<StudioProjectKind, string> = {
  dictionary: "boanerges.agentic.projectFolder.dictionary",
  topical: "boanerges.agentic.projectFolder.topical",
};

function storage(): Storage | null {
  try {
    return localStorage;
  } catch {
    return null;
  }
}

export function loadStudioFolder(kind: StudioProjectKind): string {
  const store = storage();
  if (!store) {
    return "";
  }
  return store.getItem(keys[kind]) ?? "";
}

export function saveStudioFolder(kind: StudioProjectKind, folder: string): void {
  const store = storage();
  if (!store) {
    return;
  }
  if (folder) {
    store.setItem(keys[kind], folder);
  } else {
    store.removeItem(keys[kind]);
  }
}

export function clearStudioFolder(kind: StudioProjectKind): void {
  saveStudioFolder(kind, "");
}
