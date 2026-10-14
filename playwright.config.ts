import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "e2e",
  fullyParallel: false,
  workers: 1,
  timeout: 60_000,
  expect: { timeout: 10_000 },
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 1 : 0,
  use: {
    baseURL: "http://127.0.0.1:5199",
    trace: "on-first-retry",
    actionTimeout: 15_000,
  },
  webServer: {
    command: "VITE_E2E=1 bunx vite --port 5199 --strictPort --host 127.0.0.1",
    url: "http://127.0.0.1:5199/",
    reuseExistingServer: false,
    timeout: 120_000,
  },
  projects: [{ name: "chromium", use: { ...devices["Desktop Chrome"] } }],
});
