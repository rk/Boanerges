export type DictionaryLookupResult = {
  found: boolean;
  word: string;
  variants?: { definitions: string[]; partOfSpeechExpanded: string }[];
  message?: string;
};

export function parseDictionaryLookup(resp: unknown, word: string): DictionaryLookupResult {
  if (!resp || typeof resp !== "object") {
    return { found: false, word, message: "Lookup failed." };
  }
  const o = resp as Record<string, unknown>;
  if (o.status === "importing") {
    return {
      found: false,
      word,
      message: typeof o.message === "string" ? o.message : "Dictionary is still being set up.",
    };
  }
  return resp as DictionaryLookupResult;
}
