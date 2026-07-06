export function normalizeDictionaryWord(input: string): string {
    const trimmed = input.trim();

    if (trimmed === '') {
        return '';
    }

    const token = trimmed.split(/\s+/u)[0] ?? '';
    const cleaned = token.replace(/^[^\p{L}\p{N}'-]+|[^\p{L}\p{N}'-]+$/gu, '');

    return cleaned;
}
