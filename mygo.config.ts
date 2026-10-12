import { defineConfig } from "mygo-cli";

export default defineConfig({
  name: "Boanerges",
  identifier: "com.rdk.boanerges",
  version: "0.5.0",
  copyright: "Copyright 2026 by Robert Kosek",
  comment: "A unique Bible study app.",
  devUrl: "http://localhost:5173",
  devCommand: "bun run dev:web",
  buildCommand: "bun run build:web",
  frontendDist: "dist",
  bindings: "src/mygo.ts",
  out: "build",
  resources: ["extras"],
});
