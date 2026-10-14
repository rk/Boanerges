import { AppChrome } from "../mygo";
import type { ReadabilitySettings, StudySettings } from "../mygo";
import { readability } from "./readability.svelte";
import { study } from "./study.svelte";

export type PrintMode = "include-user-work" | "blank-writing";

function studyPayload(): StudySettings {
  return {
    columnCount: study.columnCount,
    columns: study.columns,
    bookId: study.bookId,
    chapter: study.chapter,
    translationId: study.translationId,
    translationBId: study.translationBId,
    translationCId: study.translationCId,
    verseListShowContent: study.verseListShowContent,
    verseListActiveId: study.verseListActiveId,
  };
}

function readabilityPayload(): ReadabilitySettings {
  return {
    fontSize: readability.fontSize,
    lineHeight: readability.lineHeight,
    theme: readability.theme,
    fontFamily: readability.fontFamily,
    justifyText: readability.justifyText,
  };
}

export async function printStudy(mode: PrintMode): Promise<void> {
  await AppChrome.printStudy(
    studyPayload(),
    readabilityPayload(),
    mode === "include-user-work",
  );
}

export async function saveStudyPrintHtml(includeUserWork: boolean): Promise<string> {
  return AppChrome.writeStudyPrintFile(
    studyPayload(),
    readabilityPayload(),
    includeUserWork,
    "study-print.html",
  );
}
