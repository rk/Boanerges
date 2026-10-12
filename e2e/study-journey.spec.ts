import { expect, test } from "@playwright/test";

test.beforeEach(async ({ page }) => {
  await page.goto("/");
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
  await expect(page.getByText("John")).toBeVisible();
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
  await expect(page.getByTestId("comparison-pane")).toHaveCount(2);
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
  await page.waitForTimeout(500);
  await page.reload();
  await expect(page.getByTestId("study-screen")).toBeVisible({ timeout: 15_000 });
  await expect(page.getByTestId("notes-editor")).toHaveValue("e2e note", { timeout: 10_000 });
});

test("print dialog opens", async ({ page }) => {
  await page.evaluate(() => window.__dispatchStudyMenu?.("study.print"));
  await expect(page.getByTestId("print-options")).toBeVisible();
  await page.getByRole("button", { name: "Print…" }).click();
  await expect(page.getByTestId("print-options")).toBeHidden();
});
