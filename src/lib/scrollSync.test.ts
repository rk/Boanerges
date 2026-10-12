import { describe, expect, it } from "vitest";
import { syncScrollFromSource } from "./scrollSync";

describe("syncScrollFromSource", () => {
  it("mirrors scroll ratio to targets", () => {
    const source = document.createElement("div");
    Object.defineProperty(source, "scrollHeight", { value: 200 });
    Object.defineProperty(source, "clientHeight", { value: 100 });
    source.scrollTop = 50;

    const target = document.createElement("div");
    Object.defineProperty(target, "scrollHeight", { value: 400 });
    Object.defineProperty(target, "clientHeight", { value: 100 });

    syncScrollFromSource(source, [target]);
    expect(target.scrollTop).toBe(150);
  });
});
