import { describe, expect, it } from "vitest";
import { parseVerseHtml } from "./parseVerseHtml";

describe("parseVerseHtml", () => {
  it("allows legacy markup tags", () => {
    const nodes = parseVerseHtml(`In the <em>beginning</em> <span class="words-of-jesus">God</span> created.`);
    expect(nodes.some((n) => n.type === "tag" && n.name === "em")).toBe(true);
    expect(nodes.some((n) => n.type === "tag" && n.name === "span" && n.className === "words-of-jesus")).toBe(true);
  });

  it("strips unknown tags but keeps text", () => {
    const nodes = parseVerseHtml(`<script>x</script>plain`);
    expect(nodes).toEqual([{ type: "text", value: "plain" }]);
  });
});
