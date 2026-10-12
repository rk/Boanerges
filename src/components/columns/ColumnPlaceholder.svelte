<script lang="ts">
  let { columnType = "column" }: { columnType?: string; translationId?: string; slotIndex?: number; primary?: boolean } = $props();
</script>

<div class="p-4 text-sm text-base-content/60">{columnType} column</div>
