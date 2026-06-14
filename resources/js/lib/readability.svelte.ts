import { http } from '@inertiajs/svelte';

import { updateReadability as updateReadabilitySettings } from '@/actions/App/Http/Controllers/SettingsController';
import type { ReadabilitySettings, ReaderFontFamily, ReaderTheme } from '@/lib/types/readability';

const fontStacks: Record<ReaderFontFamily, string> = {
    'sans-serif':
        'Instrument Sans, ui-sans-serif, system-ui, sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji"',
    serif: "'Iowan Old Style', 'Palatino Linotype', Palatino, Georgia, serif",
};

const defaults: ReadabilitySettings = {
    fontSize: 18,
    lineHeight: 1.7,
    theme: 'light',
    fontFamily: 'serif',
};

export type { ReaderFontFamily, ReaderTheme, ReadabilitySettings };

export const readability = $state({ ...defaults });

let hydrated = false;
let persistTimeout: ReturnType<typeof setTimeout> | null = null;

function applyTheme(theme: ReaderTheme): void {
    if (typeof document === 'undefined') {
        return;
    }

    document.documentElement.setAttribute('data-theme', theme);
}

export function hydrateReadability(settings: ReadabilitySettings): void {
    if (hydrated) {
        return;
    }

    readability.fontSize = settings.fontSize;
    readability.lineHeight = settings.lineHeight;
    readability.theme = settings.theme;
    readability.fontFamily = settings.fontFamily;
    applyTheme(readability.theme);
    hydrated = true;
}

export function getReaderFontStack(family: ReaderFontFamily = readability.fontFamily): string {
    return fontStacks[family];
}

export function getReaderStyle(): string {
    return [
        `--reader-font-size: ${readability.fontSize}px`,
        `--reader-line-height: ${readability.lineHeight}`,
        `--reader-font-family: ${fontStacks[readability.fontFamily]}`,
    ].join('; ');
}

function schedulePersist(): void {
    if (! hydrated) {
        return;
    }

    if (persistTimeout) {
        clearTimeout(persistTimeout);
    }

    persistTimeout = setTimeout(() => {
        http.patch(updateReadabilitySettings.url(), {
            fontSize: readability.fontSize,
            lineHeight: readability.lineHeight,
            theme: readability.theme,
            fontFamily: readability.fontFamily,
        });
    }, 300);
}

export function setFontSize(value: number): void {
    readability.fontSize = value;
    schedulePersist();
}

export function setLineHeight(value: number): void {
    readability.lineHeight = value;
    schedulePersist();
}

export function setTheme(value: ReaderTheme): void {
    readability.theme = value;
    applyTheme(value);
    schedulePersist();
}

export function setFontFamily(value: ReaderFontFamily): void {
    readability.fontFamily = value;
    schedulePersist();
}

if (typeof document !== 'undefined') {
    applyTheme(readability.theme);
}
