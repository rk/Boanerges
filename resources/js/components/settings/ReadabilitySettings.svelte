<script lang="ts">
    import {
        getReaderFontStack,
        readability,
        setFontFamily,
        setFontSize,
        setLineHeight,
        setTheme,
        setJustifyText
    } from '@/lib/readability.svelte.ts';

    import type {ReaderFontFamily, ReaderTheme} from '@/lib/readability.svelte.ts';

    let { onclose }: { onclose: () => void } = $props();

    let dialog: HTMLDialogElement;

    const fontOptions: { id: ReaderFontFamily; label: string }[] = [
        { id: 'sans-serif', label: 'Sans-serif' },
        { id: 'serif', label: 'Serif' },
    ];

    const themeOptions: { id: ReaderTheme; label: string }[] = [
        { id: 'light', label: 'Light' },
        { id: 'dark', label: 'Dark' },
        { id: 'sepia', label: 'Sepia' },
    ];

    $effect(() => {
        dialog?.showModal();
    });

    function handleClose(): void {
        dialog?.close();
        onclose();
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box">
        <h2 class="mb-4 text-lg font-bold">Edit Settings</h2>

        <div class="space-y-6">
            <label class="input w-full">
                <b class="label">Font Size</b>
                <span class="badge badge-soft badge-primary">{readability.fontSize}px</span>
                <input
                    type="range"
                    min="14"
                    max="24"
                    step="1"
                    class="range range-sm"
                    value={readability.fontSize}
                    oninput={(event) => setFontSize(Number(event.currentTarget.value))}
                />
            </label>

            <label class="input w-full">
                <b class="label">Line Height</b>
                <span class="badge badge-soft badge-primary">{readability.lineHeight.toFixed(1)}</span>
                <input
                    type="range"
                    min="1.4"
                    max="2"
                    step="0.1"
                    class="range range-sm"
                    value={readability.lineHeight}
                    oninput={(event) => setLineHeight(Number(event.currentTarget.value))}
                />
            </label>

            <label class="input">
                <b class="label">Justify Text</b>
                <input type="checkbox" bind:checked={readability.justifyText} class="toggle" onchange={(event) => setJustifyText(event.currentTarget.checked)} />
            </label>

            <fieldset class="fieldset">
                <legend class="fieldset-legend">Font family</legend>
                <div class="join grid w-full grid-cols-2">
                    {#each fontOptions as option (option.id)}
                        <button
                            type="button"
                            class="btn join-item h-auto min-h-12 flex-col py-2"
                            class:btn-primary={readability.fontFamily === option.id}
                            style:font-family={getReaderFontStack(option.id)}
                            onclick={() => setFontFamily(option.id)}
                        >
                            {option.label}
                        </button>
                    {/each}
                </div>
                <p class="label text-xs">Preview each option in its own typeface.</p>
            </fieldset>

            <fieldset class="fieldset">
                <legend class="fieldset-legend">Theme</legend>
                <div class="join grid w-full grid-cols-3">
                    {#each themeOptions as option (option.id)}
                        <button
                            type="button"
                            class="btn join-item"
                            class:btn-primary={readability.theme === option.id}
                            onclick={() => setTheme(option.id)}
                        >
                            {option.label}
                        </button>
                    {/each}
                </div>
                <p class="label text-xs">Theme adjusts contrast and background for reading.</p>
            </fieldset>
        </div>

        <div class="modal-action">
            <form method="dialog">
                <button class="btn">Close</button>
            </form>
        </div>
    </div>

    <form method="dialog" class="modal-backdrop">
        <button aria-label="Close readability settings">close</button>
    </form>
</dialog>
