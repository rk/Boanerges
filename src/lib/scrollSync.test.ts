import { describe, expect, it } from "vitest";
import { syncLinkedScrollPanes, syncScrollFromSource } from "./scrollSync";

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

  it("syncs other panes when secondary is source", () => {
    const primary = document.createElement("div");
    Object.defineProperty(primary, "scrollHeight", { value: 200 });
    Object.defineProperty(primary, "clientHeight", { value: 100 });

    const secondary = document.createElement("div");
    Object.defineProperty(secondary, "scrollHeight", { value: 400 });
    Object.defineProperty(secondary, "clientHeight", { value: 100 });
    secondary.scrollTop = 150;

    syncLinkedScrollPanes(secondary, [primary, secondary]);
    expect(primary.scrollTop).toBe(50);
  });
});
