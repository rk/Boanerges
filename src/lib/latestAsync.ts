/** ponytail: monotonic generation for "latest async completion wins"; upgrade path — Promise.withResolvers + AbortSignal if browsers/target allow. */
export function createLatestGen() {
  let gen = 0;
  return {
    begin(): number {
      return ++gen;
    },
    invalidate(): void {
      gen++;
    },
    isStale(id: number): boolean {
      return id !== gen;
    },
  };
}
