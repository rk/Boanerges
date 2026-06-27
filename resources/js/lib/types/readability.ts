export type ReaderFontFamily = 'sans-serif' | 'serif';

export type ReaderTheme = 'light' | 'dark' | 'sepia';

export type ReadabilitySettings = {
    fontSize: number;
    lineHeight: number;
    theme: ReaderTheme;
    fontFamily: ReaderFontFamily;
    justifyText: boolean;
};
