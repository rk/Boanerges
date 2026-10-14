/** Apply daisyUI theme on <html> (light, dark, or system via --prefersdark). */
export function applyReaderTheme(theme: string): () => void {
  if (typeof document === "undefined") {
    return () => {};
  }

  const root = document.documentElement;
  const mode = theme === "light" || theme === "dark" ? theme : "auto";

  if (mode === "auto") {
    root.removeAttribute("data-theme");
    const mq = window.matchMedia("(prefers-color-scheme: dark)");
    const onChange = () => root.removeAttribute("data-theme");
    mq.addEventListener("change", onChange);
    return () => mq.removeEventListener("change", onChange);
  }

  root.setAttribute("data-theme", mode);
  return () => {};
}
