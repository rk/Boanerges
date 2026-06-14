<script lang="ts">
    let {
        title,
        translationAbbrev,
    }: {
        title: string;
        translationAbbrev?: string;
    } = $props();

    const heading = $derived(
        translationAbbrev ? `${title} (${translationAbbrev})` : title,
    );
</script>

<h1 class="reader-prose mb-4 text-center text-2xl font-semibold">{heading}</h1>
