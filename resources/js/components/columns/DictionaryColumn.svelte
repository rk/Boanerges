<script lang="ts">
    import BookOpen from '@lucide/svelte/icons/book-open';
    import ColumnHeader from '@/components/layout/ColumnHeader.svelte';
    import {
        dictionary,
        loadDictionaryEntry,
        scheduleDictionaryLookup,
        scheduleDictionarySuggest,
    } from '@/lib/dictionary.svelte.ts';
    import { normalizeDictionaryWord } from '@/lib/normalizeDictionaryWord';

    let {
        slotIndex,
    }: {
        slotIndex: number;
    } = $props();

    let wordInput = $state('');
    let menuOpen = $state(false);
    let highlightedIndex = $state(-1);

    $effect(() => {
        if (dictionary.activeWord !== null) {
            wordInput = dictionary.activeWord;
            dictionary.activeWord = null;
        }
    });

    $effect(() => {
        scheduleDictionaryLookup(wordInput);
    });

    $effect(() => {
        scheduleDictionarySuggest(wordInput);
        highlightedIndex = -1;
    });

    function selectSuggestion(word: string): void {
        wordInput = word;
        menuOpen = false;
        void loadDictionaryEntry(word);
    }

    function handleInputKeydown(event: KeyboardEvent): void {
        if (!menuOpen || dictionary.suggestions.length === 0) {
            if (event.key === 'Enter') {
                const word = normalizeDictionaryWord(wordInput);

                if (word !== '') {
                    void loadDictionaryEntry(word);
                }
            }

            return;
        }

        if (event.key === 'ArrowDown') {
            event.preventDefault();
            highlightedIndex = Math.min(
                highlightedIndex + 1,
                dictionary.suggestions.length - 1,
            );
        } else if (event.key === 'ArrowUp') {
            event.preventDefault();
            highlightedIndex = Math.max(highlightedIndex - 1, 0);
        } else if (event.key === 'Enter' && highlightedIndex >= 0) {
            event.preventDefault();
            selectSuggestion(dictionary.suggestions[highlightedIndex]);
        } else if (event.key === 'Escape') {
            menuOpen = false;
        }
    }
</script>

<div class="flex h-full min-h-0 min-w-0 flex-col">
    <ColumnHeader contentType="dictionary" {slotIndex} showViewSelector>
        <div
            class="dropdown min-w-0 flex-1"
            class:dropdown-open={menuOpen &&
                dictionary.suggestions.length > 0 &&
                wordInput.trim().length >= 2}
        >
            <label class="input input-bordered input-sm min-w-0 w-full">
                <BookOpen size={14} aria-hidden="true" />
                <input
                    type="search"
                    class="grow py-2"
                    placeholder="Look up a word…"
                    bind:value={wordInput}
                    aria-label="Dictionary word lookup"
                    aria-controls="dictionary-suggestions"
                    aria-autocomplete="list"
                    onfocus={() => (menuOpen = true)}
                    onblur={() => {
                        setTimeout(() => {
                            menuOpen = false;
                        }, 150);
                    }}
                    onkeydown={handleInputKeydown}
                />
            </label>
            <ul
                id="dictionary-suggestions"
                class="dropdown-content menu bg-base-100 rounded-box z-50 mt-1 max-h-64 w-full overflow-y-auto border border-base-300 p-1 shadow-lg"
                role="listbox"
            >
                {#each dictionary.suggestions as suggestion, index (suggestion)}
                    <li
                        role="option"
                        aria-selected={index === highlightedIndex}
                    >
                        <button
                            type="button"
                            class:menu-active={index === highlightedIndex}
                            onclick={() => selectSuggestion(suggestion)}
                        >
                            {suggestion}
                        </button>
                    </li>
                {/each}
            </ul>
        </div>
    </ColumnHeader>

    <div class="min-h-0 flex-1 overflow-y-auto px-2 py-2">
        {#if dictionary.importing}
            <p class="text-base-content/60 px-2 py-4 text-sm">
                Setting up Webster's 1828 dictionary…
            </p>
        {:else if dictionary.loading}
            <div class="flex justify-center py-8">
                <span class="loading loading-spinner loading-md text-primary"
                ></span>
            </div>
        {:else if wordInput.trim() === ''}
            <p class="text-base-content/60 px-2 py-4 text-sm">
                Type a word or right-click a verse to define.
            </p>
        {:else if dictionary.result?.found === false}
            <p class="text-base-content/60 px-2 py-4 text-sm">
                {dictionary.result.message ??
                    'No definition found for that word.'}
            </p>
        {:else if dictionary.result?.found}
            <article class="px-2 py-2">
                <header class="mb-4">
                    <h2 class="text-xl font-semibold">
                        {dictionary.result.word}
                    </h2>
                    <p class="text-base-content/60 mt-1 text-xs">
                        Webster's 1828 Dictionary
                    </p>
                </header>

                <div class="space-y-4">
                    {#each dictionary.result.variants ?? [] as variant, variantIndex (`${variant.partOfSpeech}-${variantIndex}`)}
                        <section class="rounded-lg border border-base-300 p-3">
                            {#if variant.partOfSpeechExpanded}
                                <span class="badge badge-ghost badge-sm mb-2">
                                    {variant.partOfSpeechExpanded}
                                </span>
                            {/if}
                            <ol class="list-decimal space-y-2 ps-5 text-sm">
                                {#each variant.definitions as definition, definitionIndex (definitionIndex)}
                                    <li class="whitespace-pre-line">
                                        {definition}
                                    </li>
                                {/each}
                            </ol>
                        </section>
                    {/each}
                </div>
            </article>
        {/if}
    </div>
</div>
