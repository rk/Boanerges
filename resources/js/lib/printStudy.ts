import { index as printersRoute, store as printStudyRoute } from '@/actions/App/Http/Controllers/StudyPrintController';
import { study } from '@/lib/study.svelte.ts';
import { normalizeColumns } from '@/lib/studyLayout';

export type PrintMode = 'include-user-work' | 'blank-writing';

export const PRINT_TO_PDF = '__pdf__';

export type StudyPrinter = {
    name: string;
    displayName: string;
    description: string;
};

type PrintStudyResponse = {
    path?: string;
};

async function jsonFetch<T>(url: string, init?: RequestInit): Promise<T> {
    const response = await fetch(url, {
        ...init,
        headers: {
            Accept: 'application/json',
            'X-Requested-With': 'XMLHttpRequest',
            'Content-Type': 'application/json',
            ...(init?.headers ?? {}),
        },
    });

    if (! response.ok) {
        const body = (await response.json().catch(() => null)) as { message?: string } | null;

        throw new Error(body?.message ?? `Request failed: ${response.status}`);
    }

    if (response.status === 204) {
        return undefined as T;
    }

    return response.json() as Promise<T>;
}

export async function fetchStudyPrinters(): Promise<StudyPrinter[]> {
    const data = await jsonFetch<{ printers: StudyPrinter[] }>(printersRoute.url());

    return data.printers;
}

export async function printStudy(mode: PrintMode, printerName = ''): Promise<string | null> {
    const data = await jsonFetch<PrintStudyResponse | undefined>(printStudyRoute.url(), {
        method: 'POST',
        body: JSON.stringify({
            includeUserWork: mode === 'include-user-work',
            printerName: printerName || null,
            columnCount: study.columnCount,
            columns: normalizeColumns(study.columnCount, study.columns),
            bookId: study.bookId,
            chapter: study.chapter,
            translationId: study.translationId,
            translationBId: study.translationBId,
            translationCId: study.translationCId,
        }),
    });

    return data?.path ?? null;
}
