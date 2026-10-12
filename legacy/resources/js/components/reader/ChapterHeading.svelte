<script lang="ts">
    let {
        title,
        translationAbbrev,
        compact,
    }: {
        title: string;
        translationAbbrev?: string;
        compact?: boolean;
    } = $props();

    const heading = $derived(
        translationAbbrev ? `${title} (${translationAbbrev})` : title,
    );
</script>

<h1
    class={[
        'reader-prose text-center font-semibold',
        compact ? 'm-0 text-lg' : 'mb-4 text-2xl',
    ]}
>
    {heading}
</h1>
