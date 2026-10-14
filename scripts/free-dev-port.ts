/** Kill anything listening on the Vite port so `mygo dev` can start a fresh dev server. */
const port = 5173;

function pidsOnPort(): number[] {
  if (process.platform === "win32") {
    return [];
  }
  try {
    const out = Bun.spawnSync(["lsof", "-ti", `:${port}`], { stdout: "pipe" }).stdout.toString().trim();
    if (!out) {
      return [];
    }
    return out
      .split("\n")
      .map((s) => Number.parseInt(s, 10))
      .filter((n) => Number.isFinite(n));
  } catch {
    return [];
  }
}

for (const pid of pidsOnPort()) {
  try {
    process.kill(pid, "SIGTERM");
  } catch {
    /* already gone */
  }
}
