const STORAGE_KEY = 'boanerges.readability';

export type ReaderFontFamily = 'sans-serif' | 'serif';

type ReadabilitySettings = {
    fontSize: number;
    lineHeight: number;
    theme: 'light' | 'dark';
    fontFamily: ReaderFontFamily;
};

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

function loadSettings(): ReadabilitySettings {
    if (typeof window === 'undefined') {
        return defaults;
    }

    try {
        const stored = localStorage.getItem(STORAGE_KEY);

        if (! stored) {
            return defaults;
        }

        return { ...defaults, ...JSON.parse(stored) };
    } catch {
        return defaults;
    }
}

function persist(settings: ReadabilitySettings): void {
    if (typeof window === 'undefined') {
        return;
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
}

function applyTheme(theme: 'light' | 'dark'): void {
    if (typeof document === 'undefined') {
        return;
    }

    document.documentElement.setAttribute('data-theme', theme);
}

const initial = loadSettings();

export const readability = $state({
    fontSize: initial.fontSize,
    lineHeight: initial.lineHeight,
    theme: initial.theme,
    fontFamily: initial.fontFamily,
});

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

export function setFontSize(value: number): void {
    readability.fontSize = value;
    persist(readability);
}

export function setLineHeight(value: number): void {
    readability.lineHeight = value;
    persist(readability);
}

export function setTheme(value: 'light' | 'dark'): void {
    readability.theme = value;
    applyTheme(value);
    persist(readability);
}

export function setFontFamily(value: ReaderFontFamily): void {
    readability.fontFamily = value;
    persist(readability);
}

if (typeof document !== 'undefined') {
    applyTheme(readability.theme);
}
