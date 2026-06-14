<script lang="ts">
    import type { Snippet } from 'svelte';

    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import ChapterNavDivider from '@/components/reader/ChapterNavDivider.svelte';
    import type { Chapter } from '@/lib/types/bible';

    let {
        chapter,
        translationAbbrev,
        prevLabel,
        nextLabel,
        onprev,
        onnext,
        onscroll,
        children,
        scrollRef = $bindable<HTMLElement | null>(null),
    }: {
        chapter: Chapter;
        translationAbbrev?: string;
        prevLabel?: string;
        nextLabel?: string;
        onprev?: () => void;
        onnext?: () => void;
        onscroll?: () => void;
        children: Snippet;
        scrollRef?: HTMLElement | null;
    } = $props();
</script>

<div class="flex h-full min-h-0 flex-col">
    {#if prevLabel && onprev}
        <ChapterNavDivider label={prevLabel} onclick={onprev} />
    {/if}

    <div
        bind:this={scrollRef}
        class="min-h-0 flex-1 overflow-y-auto px-6 py-4"
        onscroll={onscroll}
    >
        <ChapterHeading
            title="{chapter.book} {chapter.chapter}"
            {translationAbbrev}
        />
        {@render children()}
    </div>

    {#if nextLabel && onnext}
        <ChapterNavDivider label={nextLabel} onclick={onnext} />
    {/if}
</div>
