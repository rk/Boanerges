<script lang="ts">
  import { AppChrome } from "../mygo";

  const params = new URLSearchParams(typeof window !== "undefined" ? window.location.search : "");
  const bookId = params.get("bookId") ?? "gen";
  const chapter = Number(params.get("chapter") ?? "1");
  const verse = Number(params.get("verse") ?? "1");
  const word = params.get("word") ?? "";

  function run(action: string): void {
    void AppChrome.performVerseAction(action, bookId, chapter, verse, word);
  }
</script>

<ul class="menu bg-base-100 rounded-box h-full w-full p-2" role="menu">
  <li>
    <button type="button" role="menuitem" onclick={() => run("cross-references")}>Cross-references</button>
  </li>
  <li>
    <button type="button" role="menuitem" onclick={() => run("verse-list")}>Add to verse list</button>
  </li>
  <li>
    <button type="button" role="menuitem" onclick={() => run("copy-reference")}>Copy reference</button>
  </li>
  {#if word.length >= 2}
    <li>
      <button type="button" role="menuitem" onclick={() => run("dictionary")}>Dictionary: {word}</button>
    </li>
  {/if}
</ul>
