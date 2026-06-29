<script lang="ts">
    import Pilcrow from '@lucide/svelte/icons/pilcrow';

    import {
        applyParagraphStart,
        hydrateVerseSpan,
        normalizePasteText,
        paragraphStartAttribute,
        readVerseSpan,
    } from '@/lib/scribeDocument.ts';

    let {
        verseNumber,
        paragraphStart = false,
        oninput,
        onToggleParagraph,
    }: {
        verseNumber: number;
        paragraphStart?: boolean;
        oninput?: (text: string) => void;
        onToggleParagraph?: () => void;
    } = $props();

    let verseElement: HTMLSpanElement | undefined = $state();
    let wrapElement: HTMLSpanElement | undefined = $state();

    $effect(() => {
        if (!wrapElement) {
            return;
        }

        applyParagraphStart(wrapElement, paragraphStart);
    });

    export function setText(
        text: string,
        options: { force?: boolean } = {},
    ): void {
        if (!verseElement) {
            return;
        }

        hydrateVerseSpan(verseElement, text, options);
    }

    function handleInput(event: Event): void {
        const target = event.currentTarget as HTMLSpanElement;
        oninput?.(readVerseSpan(target));
    }

    function handlePaste(event: ClipboardEvent): void {
        event.preventDefault();

        const pasted = event.clipboardData?.getData('text/plain') ?? '';

        if (pasted === '') {
            return;
        }

        const selection = document.getSelection();

        if (!selection || selection.rangeCount === 0) {
            return;
        }

        selection.deleteFromDocument();
        selection
            .getRangeAt(0)
            .insertNode(document.createTextNode(normalizePasteText(pasted)));
        selection.collapseToEnd();

        if (verseElement) {
            oninput?.(readVerseSpan(verseElement));
        }
    }

    function handleKeydown(event: KeyboardEvent): void {
        if (event.key !== 'Tab' || !verseElement) {
            return;
        }

        event.preventDefault();
        event.stopPropagation();

        const root = verseElement.closest('.scribe-document');

        if (!root) {
            return;
        }

        const spans = [...root.querySelectorAll<HTMLElement>('[data-verse]')];
        const index = spans.indexOf(verseElement);
        const nextIndex = event.shiftKey ? index - 1 : index + 1;
        const next = spans[nextIndex];

        if (next) {
            next.focus();
        }
    }
</script>

<span
    bind:this={wrapElement}
    class="scribe-verse-wrap"
    data-paragraph-start={paragraphStartAttribute(paragraphStart)}
>
    <span contenteditable="false" class="scribe-verse-chrome">
        <button
            type="button"
            class="btn btn-ghost btn-xs px-0"
            class:opacity-30={!paragraphStart}
            class:opacity-100={paragraphStart}
            disabled={verseNumber === 1}
            tabindex="-1"
            aria-label={paragraphStart
                ? 'Remove paragraph break before this verse'
                : 'Start a new paragraph before this verse'}
            onclick={() => onToggleParagraph?.()}
        >
            <Pilcrow size={14} aria-hidden="true" />
        </button>
    </span>
    <span
        bind:this={verseElement}
        class="scribe-verse"
        contenteditable="plaintext-only"
        data-verse={verseNumber}
        role="textbox"
        tabindex="0"
        aria-multiline="true"
        aria-label="Verse {verseNumber}"
        spellcheck="true"
        oninput={handleInput}
        onpaste={handlePaste}
        onkeydown={handleKeydown}
    ></span>
</span>
