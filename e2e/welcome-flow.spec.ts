import { expect, test } from "@playwright/test";

test("welcome progress then study", async ({ page }) => {
  await page.goto("/?e2eWelcome=1");
  await expect(page.getByTestId("welcome-screen")).toBeVisible();
  await expect(page.getByRole("progressbar")).toBeVisible();
  await expect(page.getByTestId("study-screen")).toBeVisible({ timeout: 10_000 });
  await expect(page.getByTestId("welcome-screen")).toHaveCount(0);
});
