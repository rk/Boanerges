<script lang="ts">
    import ScribeEditor from '@/components/scribe/ScribeEditor.svelte';
    import ReaderPane from '@/components/reader/ReaderPane.svelte';
    import VerseText from '@/components/reader/VerseText.svelte';
    import { translations } from '@/lib/mock/chapter';
    import { getReaderStyle } from '@/lib/readability.svelte.ts';
    import { getCurrentChapter, study } from '@/lib/study.svelte.ts';

    const currentChapter = $derived(getCurrentChapter());
    const readerStyle = $derived(getReaderStyle());

    const translationA = $derived(translations.find((item) => item.id === study.translationId));
    const translationB = $derived(translations.find((item) => item.id === study.translationBId));
</script>

<div class="grid h-full min-h-0 grid-cols-3" style={readerStyle}>
    <ReaderPane chapter={currentChapter} translationAbbrev={translationA?.abbrev}>
        <VerseText verses={currentChapter.verses} />
    </ReaderPane>

    <ScribeEditor />

    <ReaderPane chapter={currentChapter} translationAbbrev={translationB?.abbrev}>
        <VerseText verses={currentChapter.verses} />
    </ReaderPane>
</div>
