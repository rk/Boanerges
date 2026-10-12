<script lang="ts">
  import { BibleService, Service } from "../mygo";

  let chapterText = $state<string>("Loading…");

  $effect(() => {
    void (async () => {
      const study = await Service.getStudy();
      const ch = await BibleService.chapter(study.translationId || "asv", study.bookId || "gen", study.chapter || 1);
      chapterText = `${ch.book} ${ch.chapter}: ${ch.verses[0]?.text ?? ""}`;
    })();
  });
</script>

<div class="p-6">
  <h1 class="text-2xl font-bold">Study</h1>
  <p class="text-base-content/70 mt-2">Shell placeholder — column layout in issue 7.</p>
  <article class="prose mt-4 max-w-none font-[family-name:var(--font-reader)]">{@html chapterText}</article>
</div>
