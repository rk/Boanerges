import {
    index as printersRoute,
    store as printStudyRoute,
} from '@/actions/App/Http/Controllers/StudyPrintController';
import { study } from '@/lib/study.svelte.ts';
import { normalizeColumns } from '@/lib/studyLayout';
import { verseListPrintPayload } from '@/lib/verseList.svelte.ts';

export type PrintMode = 'include-user-work' | 'blank-writing';

export const PRINT_TO_PDF = '__pdf__';

export const PRINT_TO_HTML = '__html__';

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

    if (!response.ok) {
        const body = (await response.json().catch(() => null)) as {
            message?: string;
        } | null;

        throw new Error(body?.message ?? `Request failed: ${response.status}`);
    }

    if (response.status === 204) {
        return undefined as T;
    }

    return response.json() as Promise<T>;
}

export async function fetchStudyPrinters(): Promise<StudyPrinter[]> {
    const data = await jsonFetch<{ printers: StudyPrinter[] }>(
        printersRoute.url(),
    );

    return data.printers;
}

export async function printStudy(
    mode: PrintMode,
    printerName = '',
): Promise<string | null> {
    const columns = normalizeColumns(study.columnCount, study.columns);
    const payload: Record<string, unknown> = {
        includeUserWork: mode === 'include-user-work',
        printerName: printerName || null,
        columnCount: study.columnCount,
        columns,
        bookId: study.bookId,
        chapter: study.chapter,
        translationId: study.translationId,
        translationBId: study.translationBId,
        translationCId: study.translationCId,
    };

    if (columns.includes('verse-list')) {
        payload.verseList = verseListPrintPayload();
    }

    const data = await jsonFetch<PrintStudyResponse | undefined>(
        printStudyRoute.url(),
        {
            method: 'POST',
            body: JSON.stringify(payload),
        },
    );

    return data?.path ?? null;
}
