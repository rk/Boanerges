/** Match legacy normalizeDictionaryWord (first token, strip edge punctuation). */
export function normalizeDictionaryWord(input: string): string {
  const trimmed = input.trim();
  if (trimmed === "") {
    return "";
  }
  const token = trimmed.split(/\s+/u)[0] ?? "";
  return token.replace(/^[^\p{L}\p{N}'-]+|[^\p{L}\p{N}'-]+$/gu, "");
}
