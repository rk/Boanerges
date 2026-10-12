export type ThemeGroup = 'basics' | 'light' | 'dark';

export type ThemeOption = {
    id: ReaderTheme;
    label: string;
    group: ThemeGroup;
};

export type ReaderTheme =
    | 'auto'
    | 'light'
    | 'autumn'
    | 'bumblebee'
    | 'cupcake'
    | 'emerald'
    | 'garden'
    | 'caramellatte'
    | 'lemonade'
    | 'nord'
    | 'pastel'
    | 'silk'
    | 'valentine'
    | 'winter'
    | 'dark'
    | 'coffee'
    | 'dracula'
    | 'forest'
    | 'luxury'
    | 'night'
    | 'sunset'
    | 'synthwave';

export const basicThemeOptions: ThemeOption[] = [
    { id: 'auto', label: 'Auto', group: 'basics' },
    { id: 'light', label: 'Light', group: 'basics' },
    { id: 'dark', label: 'Dark', group: 'basics' },
];

export const lightThemeOptions: ThemeOption[] = [
    { id: 'autumn', label: 'Autumn', group: 'light' },
    { id: 'bumblebee', label: 'Bumblebee', group: 'light' },
    { id: 'cupcake', label: 'Cupcake', group: 'light' },
    { id: 'emerald', label: 'Emerald', group: 'light' },
    { id: 'garden', label: 'Garden', group: 'light' },
    { id: 'caramellatte', label: 'Latte', group: 'light' },
    { id: 'lemonade', label: 'Lemonade', group: 'light' },
    { id: 'nord', label: 'Nord', group: 'light' },
    { id: 'pastel', label: 'Pastel', group: 'light' },
    { id: 'silk', label: 'Silk', group: 'light' },
    { id: 'valentine', label: 'Valentine', group: 'light' },
    { id: 'winter', label: 'Winter', group: 'light' },
];

export const darkThemeOptions: ThemeOption[] = [
    { id: 'coffee', label: 'Coffee', group: 'dark' },
    { id: 'dracula', label: 'Dracula', group: 'dark' },
    { id: 'forest', label: 'Forest', group: 'dark' },
    { id: 'luxury', label: 'Luxury', group: 'dark' },
    { id: 'night', label: 'Night', group: 'dark' },
    { id: 'sunset', label: 'Sunset', group: 'dark' },
    { id: 'synthwave', label: 'Synthwave', group: 'dark' },
];

export function normalizeReaderTheme(theme: string): ReaderTheme {
    if (theme === 'sepia') {
        return 'caramellatte';
    }

    return theme as ReaderTheme;
}
