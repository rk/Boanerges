<script lang="ts">
    import {
        getReaderFontStack,
        readability,
        setFontFamily,
        setFontSize,
        setLineHeight,
        setTheme,
        setJustifyText,
    } from '@/lib/readability.svelte.ts';

    import type { ReaderFontFamily } from '@/lib/readability.svelte.ts';
    import type { ReaderTheme } from '@/lib/themes';
    import {
        basicThemeOptions,
        darkThemeOptions,
        lightThemeOptions,
    } from '@/lib/themes';

    let { onclose }: { onclose: () => void } = $props();

    let dialog: HTMLDialogElement;

    const fontOptions: { id: ReaderFontFamily; label: string }[] = [
        { id: 'sans-serif', label: 'Sans-serif' },
        { id: 'serif', label: 'Serif' },
    ];

    $effect(() => {
        dialog?.showModal();
    });

    function handleClose(): void {
        dialog?.close();
        onclose();
    }

    function handleThemeChange(event: Event): void {
        setTheme(
            (event.currentTarget as HTMLSelectElement).value as ReaderTheme,
        );
    }
</script>

<dialog bind:this={dialog} class="modal" onclose={handleClose}>
    <div class="modal-box">
        <h2 class="mb-4 text-lg font-bold">Edit Settings</h2>

        <div class="space-y-6">
            <label class="input w-full">
                <b class="label">Font Size</b>
                <span class="badge badge-soft badge-primary"
                    >{readability.fontSize}px</span
                >
                <input
                    type="range"
                    min="14"
                    max="24"
                    step="1"
                    class="range range-sm"
                    value={readability.fontSize}
                    oninput={(event) =>
                        setFontSize(Number(event.currentTarget.value))}
                />
            </label>

            <label class="input w-full">
                <b class="label">Line Height</b>
                <span class="badge badge-soft badge-primary"
                    >{readability.lineHeight.toFixed(1)}</span
                >
                <input
                    type="range"
                    min="1.4"
                    max="2"
                    step="0.1"
                    class="range range-sm"
                    value={readability.lineHeight}
                    oninput={(event) =>
                        setLineHeight(Number(event.currentTarget.value))}
                />
            </label>

            <label class="label">
                <b>Justify Text</b>
                <input
                    type="checkbox"
                    bind:checked={readability.justifyText}
                    class="toggle"
                    onchange={(event) =>
                        setJustifyText(event.currentTarget.checked)}
                />
            </label>

            <fieldset class="fieldset">
                <legend class="fieldset-legend">Font family</legend>
                <div class="join grid w-full grid-cols-2">
                    {#each fontOptions as option (option.id)}
                        <button
                            type="button"
                            class="btn join-item h-auto min-h-12 flex-col py-2"
                            class:btn-primary={readability.fontFamily ===
                                option.id}
                            style:font-family={getReaderFontStack(option.id)}
                            onclick={() => setFontFamily(option.id)}
                        >
                            {option.label}
                        </button>
                    {/each}
                </div>
                <p class="label text-xs">
                    Preview each option in its own typeface.
                </p>
            </fieldset>

            <fieldset class="fieldset">
                <legend class="fieldset-legend">Theme</legend>
                <select
                    class="select w-full"
                    value={readability.theme}
                    onchange={handleThemeChange}
                >
                    <optgroup label="The Basics">
                        {#each basicThemeOptions as option (option.id)}
                            <option value={option.id}>{option.label}</option>
                        {/each}
                    </optgroup>
                    <optgroup label="Light">
                        {#each lightThemeOptions as option (option.id)}
                            <option value={option.id}>{option.label}</option>
                        {/each}
                    </optgroup>
                    <optgroup label="Dark">
                        {#each darkThemeOptions as option (option.id)}
                            <option value={option.id}>{option.label}</option>
                        {/each}
                    </optgroup>
                </select>
                <p class="label text-xs">
                    Auto follows your system light or dark appearance.
                </p>
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
