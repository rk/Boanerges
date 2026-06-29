<script lang="ts">
    import AppHead from '@/components/AppHead.svelte';
    import ColumnLayout from '@/components/layout/ColumnLayout.svelte';
    import PrintFab from '@/components/print/PrintFab.svelte';
    import AppLayout from '@/layouts/AppLayout.svelte';
</script>

<AppHead title="Study" />

<AppLayout>
    <ColumnLayout />
    <PrintFab />
</AppLayout>
