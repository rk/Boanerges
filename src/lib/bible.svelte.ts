import { SvelteMap } from "svelte/reactivity";
import { BibleService, type Book } from "../mygo";

export const bible = $state({
  books: [] as Book[],
  booksLoading: false,
});

const booksByTranslation = new SvelteMap<string, Book[]>();
const booksInflight = new SvelteMap<string, Promise<Book[]>>();

export function getBooksForTranslation(translationId: string): Book[] | undefined {
  return booksByTranslation.get(translationId);
}

export async function fetchBooksForTranslations(translationIds: string[]): Promise<void> {
  const missing = translationIds.filter((id) => !booksByTranslation.has(id));
  if (missing.length === 0) {
    return;
  }

  bible.booksLoading = true;
  try {
    await Promise.all(
      missing.map(async (id) => {
        let inflight = booksInflight.get(id);
        if (!inflight) {
          inflight = BibleService.books(id).then((resp) => resp.books);
          booksInflight.set(id, inflight);
        }
        const books = await inflight;
        booksByTranslation.set(id, books);
        booksInflight.delete(id);
      }),
    );
    const primary = translationIds[0];
    if (primary) {
      bible.books = booksByTranslation.get(primary) ?? [];
    }
  } finally {
    bible.booksLoading = false;
  }
}
