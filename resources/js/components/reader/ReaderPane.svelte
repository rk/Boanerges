<script lang="ts">
    import type { Snippet } from 'svelte';

    import ChapterHeading from '@/components/reader/ChapterHeading.svelte';
    import ChapterNavDivider from '@/components/reader/ChapterNavDivider.svelte';
    import type { Chapter, ChapterNavTarget } from '@/lib/types/bible';

    let {
        chapter,
        translationAbbrev,
        prevNav,
        nextNav,
        navLayout = 'horizontal',
        onprev,
        onnext,
        onscroll,
        children,
        scrollRef = $bindable<HTMLElement | null>(null),
    }: {
        chapter: Chapter;
        translationAbbrev?: string;
        prevNav?: ChapterNavTarget;
        nextNav?: ChapterNavTarget;
        navLayout?: 'horizontal' | 'vertical';
        onprev?: () => void;
        onnext?: () => void;
        onscroll?: () => void;
        children: Snippet;
        scrollRef?: HTMLElement | null;
    } = $props();
</script>

<div class="flex h-full min-h-0 flex-col">
    {#if prevNav && onprev}
        <ChapterNavDivider
            direction="prev"
            bookAbbrev={prevNav.bookAbbrev}
            chapter={prevNav.chapter}
            layout={navLayout}
            onclick={onprev}
        />
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

    {#if nextNav && onnext}
        <ChapterNavDivider
            direction="next"
            bookAbbrev={nextNav.bookAbbrev}
            chapter={nextNav.chapter}
            layout={navLayout}
            onclick={onnext}
        />
    {/if}
</div>
