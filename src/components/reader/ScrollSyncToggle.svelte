<script lang="ts">
  import Link2 from "@lucide/svelte/icons/link-2";
  import Unlink from "@lucide/svelte/icons/unlink";
  import { setScrollSync, study } from "../../lib/study.svelte.ts";

  function handleChange(event: Event): void {
    setScrollSync((event.currentTarget as HTMLInputElement).checked);
  }
</script>

<label class="label cursor-pointer gap-2 py-0">
  {#if study.scrollSync}
    <Link2 size={16} aria-hidden="true" />
  {:else}
    <Unlink size={16} aria-hidden="true" />
  {/if}
  <span class="label-text text-sm hidden sm:inline">Link scroll</span>
  <input
    type="checkbox"
    class="toggle toggle-sm"
    data-testid="scroll-sync-toggle"
    checked={study.scrollSync}
    onchange={handleChange}
  />
</label>
