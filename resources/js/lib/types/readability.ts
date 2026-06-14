export type ReaderFontFamily = 'sans-serif' | 'serif';

export type ReadabilitySettings = {
    fontSize: number;
    lineHeight: number;
    theme: 'light' | 'dark';
    fontFamily: ReaderFontFamily;
};
