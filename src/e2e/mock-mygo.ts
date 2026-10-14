import type {
  BooksResponse,
  CatalogResponse,
  ChapterResponse,
  CrossRefsResponse,
  Payload,
  ReadabilitySettings,
  Response,
  StudySettings,
} from "../mygo";

const genesisChapter: ChapterResponse = {
  chapter: {
    book: "Genesis",
    bookAbbrev: "Gen",
    chapter: 1,
    verses: [
      { number: 1, text: `In the beginning <em>God</em> created the heaven and the earth.`, paragraphStart: true },
      { number: 2, text: "And the earth was without form, and void.", paragraphStart: false },
    ],
  },
};

const store = {
  readability: {
    fontSize: 18,
    lineHeight: 1.6,
    theme: "light",
    fontFamily: "serif",
    justifyText: false,
  } as ReadabilitySettings,
  study: {
    columnCount: 1,
    columns: [] as string[],
    bookId: "gen",
    chapter: 1,
    translationId: "asv",
    translationBId: "asv",
    translationCId: "asv",
    verseListShowContent: true,
    verseListActiveId: null as string | null,
  } as StudySettings,
  notes: new Map<string, string>(),
  scribe: new Map<string, { verse: number; text: string; paragraphStart?: boolean }[]>(),
  verseLists: new Map<string, { id: string; title: string; entries: { bookId: string; chapter: number; verse: number }[]; updatedAt: string }>(),
  catalog: [
    {
      id: "asv",
      module: "ASV",
      name: "American Standard Version",
      abbrev: "ASV",
      installed: false,
      bundled: true,
      install_status: null as string | null,
    },
    {
      id: "kjv",
      module: "KJV",
      name: "King James Version",
      abbrev: "KJV",
      installed: false,
      bundled: false,
      install_status: null as string | null,
    },
    {
      id: "oeb",
      module: "OEB",
      name: "Open English Bible (US Spelling)",
      abbrev: "OEB",
      installed: false,
      bundled: false,
      install_status: null as string | null,
    },
  ],
};

function key(bookId: string, chapter: number) {
  return `${bookId}:${chapter}`;
}

declare global {
  interface Window {
    __dispatchStudyMenu?: (id: string) => void;
    __openTranslationManager?: () => Promise<void>;
    __e2eReset?: () => void;
  }
  // eslint-disable-next-line no-var
  var mygo: {
    call(method: string, ...args: unknown[]): Promise<unknown>;
    channel(onmessage?: (msg: unknown) => void): { postMessage: (msg: unknown) => void; close: () => void };
    on(name: string, listener: (payload: unknown) => void): () => void;
    once(name: string, listener: (payload: unknown) => void): () => void;
  };
}

export function installMygoMock(mode: string = "1"): void {
  const listeners = new Map<string, Set<(payload: unknown) => void>>();

  const emit = (name: string, payload: unknown) => {
    listeners.get(name)?.forEach((fn) => fn(payload));
  };

  const handlers: Record<string, (...args: unknown[]) => unknown> = {
    "BibleService.ASVInstallStatus": (): Payload =>
      mode === "welcome"
        ? { abbrev: "ASV", step: "creating_schema", percent: 30, install_status: "creating_schema" }
        : { abbrev: "ASV", step: "ready", percent: 100, install_status: "ready" },
    "BibleService.Books": (): BooksResponse => ({
      books: [
        { id: "gen", name: "Genesis", abbrev: "Gen", testament: "OT", chapters: 50, firstChapter: 1, lastChapter: 50 },
        { id: "mat", name: "Matthew", abbrev: "Matt", testament: "NT", chapters: 28, firstChapter: 1, lastChapter: 28 },
      ],
    }),
    "BibleService.Chapter": (_translation: unknown, book: unknown, chapter: unknown): ChapterResponse => {
      const ch = Number(chapter) || 1;
      const bookId = String(book);
      const bookName = bookId === "gen" ? "Genesis" : bookId;
      return {
        chapter: {
          book: bookName,
          bookAbbrev: "Gen",
          chapter: ch,
          verses: genesisChapter.chapter.verses.map((v) => ({ ...v })),
        },
      };
    },
    "BibleService.Search": (_query: unknown, _t: unknown, limit: unknown, offset: unknown): Response => {
      const off = offset as number;
      const hit = { bookId: "gen", chapter: 1, verse: 1, snippet: "In the <b>beginning</b>", translation: "asv" };
      const hit2 = { bookId: "gen", chapter: 1, verse: 2, snippet: "earth was <b>without</b> form", translation: "asv" };
      if (off === 0) {
        return { results: [hit], total: 2, hasMore: true };
      }
      return { results: [hit2], total: 2, hasMore: false };
    },
    "BibleService.CrossReferences": (): CrossRefsResponse => ({
      references: [{ rank: 1, bookId: "jhn", bookName: "John", chapter: 1, verse: 1, endVerse: null }],
    }),
    "Service.GetReadability": (): ReadabilitySettings => {
      const raw = sessionStorage.getItem("e2e.readability");
      if (raw) store.readability = JSON.parse(raw) as ReadabilitySettings;
      return { ...store.readability };
    },
    "Service.UpdateReadability": (settings: unknown): ReadabilitySettings => {
      const s = settings as ReadabilitySettings;
      if (s.fontSize < 14) {
        throw Object.assign(new Error("fontSize must be at least 14"), { name: "CallError" });
      }
      store.readability = { ...s };
      sessionStorage.setItem("e2e.readability", JSON.stringify(store.readability));
      return store.readability;
    },
    "Service.GetStudy": (): StudySettings => {
      const raw = sessionStorage.getItem("e2e.study");
      if (raw) store.study = JSON.parse(raw) as StudySettings;
      return { ...store.study, columns: [...store.study.columns] };
    },
    "Service.UpdateStudy": (settings: unknown): StudySettings => {
      const s = settings as StudySettings;
      store.study = { ...s, columns: [...s.columns] };
      sessionStorage.setItem("e2e.study", JSON.stringify(store.study));
      return store.study;
    },
    "TranslationService.Catalog": (): CatalogResponse => ({ translations: store.catalog.map((t) => ({ ...t })) }),
    "TranslationService.Install": (module: unknown) => {
      if (module === "ASV") {
        throw Object.assign(new Error("The ASV translation is bundled with the app and cannot be installed separately."), {
          name: "CallError",
        });
      }
      const mod = String(module);
      for (const entry of store.catalog) {
        if (entry.module === mod) {
          entry.installed = true;
          entry.install_status = "ready";
        }
      }
      return {
        translation: {
          id: mod.toLowerCase(),
          module: mod,
          name: mod,
          abbrev: mod,
          bundled: false,
          install_status: "ready",
        },
      };
    },
    "TranslationService.Uninstall": (module: unknown) => {
      if (module === "ASV") {
        throw Object.assign(new Error("cannot uninstall bundled ASV"), { name: "CallError" });
      }
    },
    "TranslationService.Installed": () => ({
      translations: [{ id: "asv", module: "ASV", name: "ASV", abbrev: "ASV", bundled: true, install_status: "ready" }],
    }),
    "DictionaryService.Suggest": (query: unknown) => ({ suggestions: (query as string) ? ["love", "lovely"] : [] }),
    "DictionaryService.Lookup": (word: unknown) =>
      (word as string).toLowerCase() === "love"
        ? { found: true, word: "LOVE", variants: [{ partOfSpeechExpanded: "noun", definitions: ["Strong affection."] }] }
        : { found: false, word, message: "Not found" },
    "DictionaryService.Reinstall": () => undefined,
    "DictionaryService.Status": () => ({ installed: true, entries: 113376 }),
    "StudyData.NotesGet": (bookId: unknown, chapter: unknown) => {
      const k = key(bookId as string, chapter as number);
      const persisted = sessionStorage.getItem(`e2e.notes.${k}`);
      if (persisted !== null) return persisted;
      return store.notes.get(k) ?? "";
    },
    "StudyData.NotesPut": (bookId: unknown, chapter: unknown, content: unknown) => {
      const k = key(bookId as string, chapter as number);
      store.notes.set(k, content as string);
      sessionStorage.setItem(`e2e.notes.${k}`, content as string);
    },
    "StudyData.ScribeGet": (bookId: unknown, chapter: unknown) => {
      const k = key(bookId as string, chapter as number);
      const raw = sessionStorage.getItem(`e2e.scribe.${k}`);
      return raw ? (JSON.parse(raw) as { verse: number; text: string }[]) : [];
    },
    "StudyData.ScribePut": (bookId: unknown, chapter: unknown, verses: unknown) => {
      const k = key(bookId as string, chapter as number);
      sessionStorage.setItem(`e2e.scribe.${k}`, JSON.stringify(verses));
    },
    "StudyData.VerseListIndex": () => [...store.verseLists.values()].map(({ id, title, updatedAt }) => ({ id, title, updatedAt })),
    "StudyData.VerseListSave": (title: unknown, entries: unknown, id: unknown) => {
      const list = { id: (id as string) || "list1", title: title as string, entries: entries as { bookId: string; chapter: number; verse: number }[], updatedAt: new Date().toISOString() };
      store.verseLists.set(list.id, list);
      return list;
    },
    "StudyData.VerseListGet": (id: unknown) => store.verseLists.get(id as string) ?? { id, title: "", entries: [], updatedAt: "" },
    "StudyData.VerseListDelete": (id: unknown) => {
      store.verseLists.delete(id as string);
    },
    "AppChrome.Print": () => undefined,
    "AppChrome.PrintToPDF": () => "/tmp/boanerges-study.pdf",
    "AppChrome.StudyPrintHTML": () => "<html><body><h1>Genesis 1</h1></body></html>",
    "AppChrome.WriteStudyPrintFile": () => "/tmp/boanerges-study-print.html",
    "AppChrome.OpenBookChapterFlyout": () => undefined,
    "AppChrome.CloseBookChapterFlyout": () => undefined,
    "AppChrome.OpenVerseActionsFlyout": () => undefined,
    "AppChrome.CloseVerseActionsFlyout": () => undefined,
    "AppChrome.PerformVerseAction": (action: unknown, bookId: unknown, chapter: unknown, verse: unknown, word: unknown) => {
      const payload = { action, bookId, chapter, verse, word };
      listeners.get("verseMenuAction")?.forEach((fn) => fn(payload));
    },
    "AppChrome.SyncMenuState": () => undefined,
    "AppChrome.OpenDevTools": () => undefined,
    "AppChrome.OpenFeedback": () => undefined,
  };

  globalThis.mygo = {
    call(method: string, ...args: unknown[]) {
      const fn = handlers[method];
      if (!fn) {
        return Promise.reject(new Error(`mock: unknown method ${method}`));
      }
      try {
        return Promise.resolve(fn(...args));
      } catch (err) {
        return Promise.reject(err);
      }
    },
    channel() {
      return { postMessage: () => {}, close: () => {} };
    },
    on(name: string, listener: (payload: unknown) => void) {
      if (!listeners.has(name)) listeners.set(name, new Set());
      listeners.get(name)!.add(listener);
      return () => listeners.get(name)?.delete(listener);
    },
    once(name: string, listener: (payload: unknown) => void) {
      const off = this.on(name, (payload: unknown) => {
        off();
        listener(payload);
      });
      return off;
    },
  };

  window.__e2eReset = () => {
    store.readability.fontSize = 18;
    store.study.columnCount = 1;
    store.study.columns = [];
    store.study.bookId = "gen";
    store.study.chapter = 1;
    store.study.verseListActiveId = null;
    store.notes.clear();
    store.scribe.clear();
    store.verseLists.clear();
    for (let i = sessionStorage.length - 1; i >= 0; i--) {
      const k = sessionStorage.key(i);
      if (k?.startsWith("e2e.")) {
        sessionStorage.removeItem(k);
      }
    }
  };

  window.__dispatchStudyMenu = (id: string) => {
    listeners.get("studyMenu")?.forEach((fn) => fn(id));
  };

  if (mode === "welcome") {
    queueMicrotask(() => {
      emit("translationInstallProgress", { abbrev: "ASV", step: "creating_schema", percent: 30, install_status: "creating_schema" });
      setTimeout(() => {
        emit("translationInstallProgress", { abbrev: "ASV", step: "ready", percent: 100, install_status: "ready" });
      }, 300);
    });
  }
}
