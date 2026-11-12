import { applyStudyTheme, type StudyAppearance, type StudyThemePole } from "./studyTheme";

/** @deprecated Use applyStudyTheme from studyTheme.ts */
export function applyReaderTheme(theme: string): () => void {
  const appearance: StudyAppearance =
    theme === "light" || theme === "dark" || theme === "auto" ? theme : "auto";
  return applyStudyTheme("monastic", appearance);
}

export function applyReadabilityTheme(pole: StudyThemePole, appearance: StudyAppearance): () => void {
  return applyStudyTheme(pole, appearance);
}
