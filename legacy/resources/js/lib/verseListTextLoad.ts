export function verseListNeedsTextLoad(
    entries: ReadonlyArray<{ text?: string }>,
    showContent: boolean,
    translationId: string,
    loading: boolean,
    loadedTranslationId: string | null,
): boolean {
    if (loading) {
        return false;
    }

    return (
        showContent &&
        entries.length > 0 &&
        (loadedTranslationId !== translationId ||
            entries.some((entry) => entry.text === undefined))
    );
}
