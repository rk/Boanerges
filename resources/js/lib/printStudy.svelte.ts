export const printUi = $state({
    optionsOpen: false,
});

export function openPrintOptions(): void {
    printUi.optionsOpen = true;
}

export function closePrintOptions(): void {
    printUi.optionsOpen = false;
}
