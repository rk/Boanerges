import { defineConfig } from "vitest/config";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  test: {
    environment: "happy-dom",
    include: ["src/**/*.test.ts", "e2e/**/*.test.ts"],
  },
});
