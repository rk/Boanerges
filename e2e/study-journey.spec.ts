import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("/");
  await page.evaluate(() => window.__e2eReset?.());
  await expect(page.getByTestId("study-screen")).toBeVisible({ timeout: 15_000 });
});

test("welcome skipped when ASV ready", async ({ page }) => {
  await expect(page.getByTestId("welcome-screen")).toHaveCount(0);
});

test("books → chapter → markup verse", async ({ page }) => {
  await expect(page.getByTestId("chapter-heading")).toHaveText("Genesis 1");
  await expect(page.getByTestId("bible-reader-primary")).toContainText("beginning");
  await expect(page.locator("[data-verse='1'] em")).toHaveCount(1);
});

test("search hit navigates location", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.search"));
  await page.getByTestId("search-input").fill("beginning");
  await page.getByTestId("search-input").press("Enter");
  await page.getByTestId("search-hit").click();
  await expect(page.getByTestId("chapter-heading")).toHaveText("Genesis 1");
});

test("cross references panel", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.cross-references"));
  await page.getByTestId("bible-reader-primary").locator("[data-verse='1']").click({ button: "right" });
  await page.getByRole("menuitem", { name: "Cross-references" }).click();
  await expect(page.getByTestId("cross-ref-hit").first()).toContainText("John");
});

test("dictionary lookup", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.dictionary"));
  await page.getByTestId("dictionary-input").fill("love");
  await page.getByTestId("dictionary-input").press("Enter");
  await expect(page.getByText("Strong affection")).toBeVisible();
});

test("columns and comparison", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.comparison"));
  await expect(page.getByTestId("comparison-column")).toBeVisible();
  await page.getByTestId("comparison-reference").fill("Genesis 1:1");
  await expect(page.getByTestId("comparison-row").first()).toBeVisible({ timeout: 5000 });
});

test("settings persist readability", async ({ page }) => {
  await page.getByRole("button", { name: "Settings" }).click();
  const slider = page.getByTestId("readability-settings").locator('input[type="range"]').first();
  await slider.fill("22");
  await page.reload();
  await page.getByRole("button", { name: "Settings" }).click();
  await expect(slider).toHaveValue("22");
});

test("translation manager bundled guard", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("file.manage-translations"));
  await expect(page.getByTestId("translation-manager")).toBeVisible();
  await page.getByRole("listitem").filter({ hasText: "American Standard Version" }).getByRole("button", { name: "Install" }).click();
  await expect(page.getByTestId("translation-manager-error")).toContainText("ASV");
});

test("notes persist per chapter", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.notes"));
  await expect(page.getByTestId("notes-editor")).toBeVisible();
  await page.getByTestId("notes-editor").fill("e2e note");
  await expect
    .poll(async () => page.evaluate(() => sessionStorage.getItem("e2e.notes.gen:1")))
    .toBe("e2e note");
  await page.reload();
  await expect(page.getByTestId("study-screen")).toBeVisible({ timeout: 15_000 });
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.notes"));
  await expect(page.getByTestId("notes-editor")).toHaveValue("e2e note", { timeout: 10_000 });
});

test("print dialog opens", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.print"));
  await expect(page.getByTestId("print-options")).toBeVisible();
  await page.getByTestId("print-include-work").click();
  await expect(page.getByTestId("print-options")).toBeHidden();
});

test("readability validation error", async ({ page }) => {
  await page.getByRole("button", { name: "Settings" }).click();
  await page.getByTestId("readability-invalid").click();
  await expect(page.getByTestId("readability-error")).toBeVisible();
});

test("chapter navigation", async ({ page }) => {
  await page.getByTestId("chapter-next").click();
  await expect(page.getByTestId("chapter-heading")).toHaveText("Genesis 2");
  await page.getByTestId("chapter-prev").click();
  await expect(page.getByTestId("chapter-heading")).toHaveText("Genesis 1");
});

test("search load more", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.search"));
  await page.getByTestId("search-input").fill("beginning");
  await page.getByTestId("search-input").press("Enter");
  await page.getByTestId("search-load-more").click();
  await expect(page.getByTestId("search-hit")).toHaveCount(2);
});

test("catalog USFM module install", async ({ page }) => {
  await page.evaluate(async () => {
    await window.__openTranslationManager?.();
  });
  const row = page.getByTestId("translation-manager").locator("li").filter({ hasText: "Open English Bible" });
  await expect(row).toBeVisible({ timeout: 10_000 });
  await row.getByRole("button", { name: "Install" }).click();
  await expect(row.getByRole("button", { name: "Remove" })).toBeVisible({ timeout: 5000 });
});

test("catalog translation install", async ({ page }) => {
  await page.evaluate(async () => {
    await window.__openTranslationManager?.();
  });
  const row = page.getByTestId("translation-manager").locator("li").filter({ hasText: "King James Version" });
  await expect(row).toBeVisible({ timeout: 10_000 });
  await row.getByRole("button", { name: "Install" }).click();
  await expect(row.getByRole("button", { name: "Remove" })).toBeVisible({ timeout: 5000 });
});

test("scribe persists", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.scribe"));
  const verse1 = page.getByRole("textbox", { name: "Verse 1" });
  await expect(verse1).toBeVisible({ timeout: 10_000 });
  await verse1.click();
  await verse1.fill("scribe e2e line");
  await verse1.blur();
  await expect
    .poll(async () => {
      const raw = await page.evaluate(() => sessionStorage.getItem("e2e.scribe.gen:1"));
      return raw?.includes("scribe e2e") ?? false;
    })
    .toBe(true);
  await page.reload();
  await expect(page.getByTestId("study-screen")).toBeVisible({ timeout: 15_000 });
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.scribe"));
  await expect(page.getByRole("textbox", { name: "Verse 1" })).toContainText("scribe e2e", { timeout: 10_000 });
});

test("verse list add entry", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("view.columns.2"));
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.verse-list"));
  await expect(page.getByTestId("verse-list-column")).toBeVisible();
  await page.getByTestId("bible-reader-primary").locator("[data-verse='1']").click();
  await page.getByRole("button", { name: "Add verse" }).click();
  await expect(page.getByText("gen 1:1")).toBeVisible();
});
