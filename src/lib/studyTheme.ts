export type StudyThemePole = "monastic" | "modern";
export type StudyAppearance = "auto" | "light" | "dark";

export function studyThemeName(pole: StudyThemePole, appearance: StudyAppearance): string {
  if (typeof window === "undefined") {
    return `study-${pole}-light`;
  }
  const mode =
    appearance === "auto"
      ? window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light"
      : appearance;
  return `study-${pole}-${mode}`;
}

/** Apply composed daisyUI theme on `<html>`. Auto follows prefers-color-scheme. */
export function applyStudyTheme(pole: StudyThemePole, appearance: StudyAppearance): () => void {
  if (typeof document === "undefined") {
    return () => {};
  }

  const root = document.documentElement;

  const apply = () => {
    root.setAttribute("data-theme", studyThemeName(pole, appearance));
  };

  apply();

  if (appearance !== "auto") {
    return () => {};
  }

  const mq = window.matchMedia("(prefers-color-scheme: dark)");
  const onChange = () => apply();
  mq.addEventListener("change", onChange);
  return () => mq.removeEventListener("change", onChange);
}
