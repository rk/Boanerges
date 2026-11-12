import { StudyData } from "../mygo";

export const notes = $state({
  saveStatus: "idle" as "idle" | "saving" | "saved",
});

let saveTimeout: ReturnType<typeof setTimeout> | null = null;
let hideTimeout: ReturnType<typeof setTimeout> | null = null;
let saveGeneration = 0;

export function scheduleNotesSave(bookId: string, chapter: number, content: string): void {
  notes.saveStatus = "saving";

  if (saveTimeout) {
    clearTimeout(saveTimeout);
  }
  if (hideTimeout) {
    clearTimeout(hideTimeout);
  }

  const generation = ++saveGeneration;

  saveTimeout = setTimeout(() => {
    void StudyData.notesPut(bookId, chapter, content)
      .then(() => {
        if (generation !== saveGeneration) {
          return;
        }

        notes.saveStatus = "saved";

        hideTimeout = setTimeout(() => {
          if (generation === saveGeneration) {
            notes.saveStatus = "idle";
          }
        }, 2000);
      })
      .catch(() => {
        if (generation === saveGeneration) {
          notes.saveStatus = "idle";
        }
      });
  }, 500);
}
