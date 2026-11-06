import { describe, expect, it } from "vitest";
import { createLatestGen } from "./latestAsync";

describe("createLatestGen", () => {
  it("treats earlier completions as stale after invalidate", async () => {
    const latest = createLatestGen();
    const id1 = latest.begin();
    const p1 = Promise.resolve("slow").then(() => ({ id: id1, val: "a" }));
    latest.invalidate();
    const id2 = latest.begin();
    expect(latest.isStale(id1)).toBe(true);
    expect(latest.isStale(id2)).toBe(false);
    const first = await p1;
    expect(latest.isStale(first.id)).toBe(true);
  });
});
