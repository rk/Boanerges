export type SnippetPart = { text: string; highlight: boolean };

const MARK_PATTERN = /<mark>(.*?)<\/mark>/gi;

export function parseHighlightSnippet(snippet: string): SnippetPart[] {
    if (snippet === '') {
        return [];
    }

    const parts: SnippetPart[] = [];
    let lastIndex = 0;

    for (const match of snippet.matchAll(MARK_PATTERN)) {
        const index = match.index ?? 0;

        if (index > lastIndex) {
            parts.push({
                text: snippet.slice(lastIndex, index),
                highlight: false,
            });
        }

        parts.push({ text: match[1] ?? '', highlight: true });
        lastIndex = index + match[0].length;
    }

    if (lastIndex < snippet.length) {
        parts.push({ text: snippet.slice(lastIndex), highlight: false });
    }

    return parts.length > 0 ? parts : [{ text: snippet, highlight: false }];
}
