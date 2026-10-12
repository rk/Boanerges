import type { ReaderTheme } from '@/lib/themes';

export type { ReaderTheme };

export type ReaderFontFamily = 'sans-serif' | 'serif';

export type ReadabilitySettings = {
    fontSize: number;
    lineHeight: number;
    theme: ReaderTheme;
    fontFamily: ReaderFontFamily;
    justifyText: boolean;
};
