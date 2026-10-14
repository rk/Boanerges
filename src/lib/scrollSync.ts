/** Sync scroll position from source to targets by scroll ratio (ponytail: linear ratio). */
export function syncScrollFromSource(source: HTMLElement, targets: HTMLElement[]): void {
  const max = source.scrollHeight - source.clientHeight;
  const ratio = max <= 0 ? 0 : source.scrollTop / max;
  for (const target of targets) {
    if (target === source) continue;
    const targetMax = target.scrollHeight - target.clientHeight;
    target.scrollTop = ratio * targetMax;
  }
}

/** Run sync without re-entering when programmatic scrollTop triggers scroll events. */
export function syncLinkedScrollPanes(source: HTMLElement, panes: HTMLElement[]): void {
  const targets = panes.filter((el) => el !== source);
  if (targets.length === 0) {
    return;
  }
  syncScrollFromSource(source, targets);
}
